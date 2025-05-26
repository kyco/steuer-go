package bmf

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

const (
	XMLSourceURL = "https://www.bmf-steuerrechner.de/javax.faces.resource/daten/xmls/Lohnsteuer2025.xml.xhtml"
)

type PAPData struct {
	XMLName    xml.Name        `xml:"PAP"`
	Name       string          `xml:"name,attr"`
	Version    string          `xml:"version,attr"`
	Variables  PAPVariables    `xml:"VARIABLES"`
	Constants  PAPConstants    `xml:"CONSTANTS"`
	Methods    PAPMethods      `xml:"METHODS"`
}

type PAPVariables struct {
	Inputs    PAPInputs    `xml:"INPUTS"`
	Outputs   PAPOutputs   `xml:"OUTPUTS"`
	Internals PAPInternals `xml:"INTERNALS"`
}

type PAPInputs struct {
	Input []InputVariable `xml:"INPUT"`
}

type PAPOutputs struct {
	Type   string           `xml:"type,attr"`
	Output []OutputVariable `xml:"OUTPUT"`
}

type PAPInternals struct {
	Internal []InternalVariable `xml:"INTERNAL"`
}

type InputVariable struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}

type OutputVariable struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}

type InternalVariable struct {
	Name    string `xml:"name,attr"`
	Type    string `xml:"type,attr"`
	Default string `xml:"default,attr"`
}

type PAPConstants struct {
	Constant []PAPConstant `xml:"CONSTANT"`
}

type PAPConstant struct {
	Name  string `xml:"name,attr"`
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type PAPMethods struct {
	Main   []PAPMethod `xml:"MAIN"`
	Method []PAPMethod `xml:"METHOD"`
}

type PAPMethod struct {
	Name  string        `xml:"name,attr"`
	Steps []interface{} `xml:",any"`
}

type TaxCalculator struct {
	XMLData      *PAPData
	InputValues  map[string]interface{}
	OutputValues map[string]interface{}
	InternalVars map[string]interface{}
	Constants    map[string]interface{}
}

func FetchTaxCalculationXML() (*PAPData, error) {
	resp, err := http.Get(XMLSourceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch XML: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("XML fetch failed with status: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read XML data: %w", err)
	}

	var papData PAPData
	if err := xml.Unmarshal(data, &papData); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}

	return &papData, nil
}

func NewTaxCalculator(papData *PAPData) *TaxCalculator {
	calculator := &TaxCalculator{
		XMLData:      papData,
		InputValues:  make(map[string]interface{}),
		OutputValues: make(map[string]interface{}),
		InternalVars: make(map[string]interface{}),
		Constants:    make(map[string]interface{}),
	}

	// Initialize default values from XML
	for _, input := range papData.Variables.Inputs.Input {
		if input.Default != "" {
			calculator.InputValues[input.Name] = parseValue(input.Type, input.Default)
		}
	}

	for _, output := range papData.Variables.Outputs.Output {
		if output.Default != "" {
			calculator.OutputValues[output.Name] = parseValue(output.Type, output.Default)
		}
	}

	for _, internal := range papData.Variables.Internals.Internal {
		if internal.Default != "" {
			calculator.InternalVars[internal.Name] = parseValue(internal.Type, internal.Default)
		}
	}

	// Initialize constants
	for _, constant := range papData.Constants.Constant {
		calculator.Constants[constant.Name] = parseValue(constant.Type, constant.Value)
	}

	return calculator
}

func (tc *TaxCalculator) SetInputValue(name string, value interface{}) {
	tc.InputValues[name] = value
}

func (tc *TaxCalculator) GetOutputValue(name string) interface{} {
	return tc.OutputValues[name]
}

func (tc *TaxCalculator) Calculate() error {
	// Reset output values and internal variables
	tc.OutputValues = make(map[string]interface{})
	tc.InternalVars = make(map[string]interface{})

	// Initialize defaults again
	for _, output := range tc.XMLData.Variables.Outputs.Output {
		if output.Default != "" {
			tc.OutputValues[output.Name] = parseValue(output.Type, output.Default)
		}
	}

	for _, internal := range tc.XMLData.Variables.Internals.Internal {
		if internal.Default != "" {
			tc.InternalVars[internal.Name] = parseValue(internal.Type, internal.Default)
		}
	}

	// Execute the main calculation method
	if len(tc.XMLData.Methods.Main) > 0 {
		// Process each step in the main method
		for _, step := range tc.XMLData.Methods.Main[0].Steps {
			if execute, ok := step.(xml.StartElement); ok {
				op := OperationType(execute.Name.Local)
				
				switch op {
				case OpExecute:
					// Call another method
					methodName := ""
					for _, attr := range execute.Attr {
						if attr.Name.Local == "method" {
							methodName = attr.Value
							break
						}
					}
					if methodName != "" {
						if err := tc.executeMethod(methodName); err != nil {
							return fmt.Errorf("error executing method %s: %w", methodName, err)
						}
					}
					
				case OpEval:
					// Evaluate an expression
					var target string
					var left, right string
					var op EvalOperator
					
					for _, attr := range execute.Attr {
						switch attr.Name.Local {
						case "target":
							target = attr.Value
						case "op":
							op = EvalOperator(attr.Value)
						case "left":
							left = attr.Value
						case "right":
							right = attr.Value
						}
					}
					
					if target != "" {
						// Evaluate the expression
						result, err := tc.evaluateExpression(left, right, op)
						if err != nil {
							return fmt.Errorf("evaluation error in main method: %w", err)
						}
						
						// Store the result in the appropriate location
						tc.setVariableValue(target, result)
					}
					
				case OpIf:
					// Process IF condition in the main method
					var left, right string
					var op ComparisonOperator
					
					for _, attr := range execute.Attr {
						switch attr.Name.Local {
						case "left":
							left = attr.Value
						case "right":
							right = attr.Value
						case "op":
							op = ComparisonOperator(attr.Value)
						}
					}
					
					// Evaluate the condition
					result, err := tc.evaluateComparison(left, right, op)
					if err != nil {
						return fmt.Errorf("if condition error in main method: %w", err)
					}
					
					// Process THEN/ELSE blocks - this is simplified and would need to be
					// implemented more completely by tracking the XML structure
					// For now, we'll just log the condition result
					_ = result
				}
			}
		}
	} else {
		return fmt.Errorf("no main method found in XML data")
	}

	// Round monetary outputs to integers (cents)
	for name, value := range tc.OutputValues {
		if rat, ok := value.(*big.Rat); ok {
			// For monetary values, round to cents
			cents := new(big.Rat).Mul(rat, big.NewRat(100, 1))
			f, _ := cents.Float64()
			tc.OutputValues[name] = int(math.Round(f))
		}
	}

	return nil
}

type OperationType string

const (
	OpExecute  OperationType = "EXECUTE"
	OpEval     OperationType = "EVAL"
	OpIf       OperationType = "IF"
	OpCompare  OperationType = "COMPARE"
	OpBausteinFinish OperationType = "BAUSTEINFINISH"
	OpThen     OperationType = "THEN"
	OpElse     OperationType = "ELSE"
)

type ComparisonOperator string

const (
	CompLT   ComparisonOperator = "LT"
	CompLE   ComparisonOperator = "LE"
	CompGT   ComparisonOperator = "GT"
	CompGE   ComparisonOperator = "GE"
	CompEQ   ComparisonOperator = "EQ"
	CompNE   ComparisonOperator = "NE"
)

type EvalOperator string

const (
	EvalAssign      EvalOperator = "="
	EvalAdd         EvalOperator = "+"
	EvalSubtract    EvalOperator = "-"
	EvalMultiply    EvalOperator = "*"
	EvalDivide      EvalOperator = "/"
	EvalIntDivide   EvalOperator = "DIV"
	EvalModulo      EvalOperator = "MOD"
)

func (tc *TaxCalculator) executeMethod(methodName string) error {
	// Find the method with the given name
	var methodToExecute *PAPMethod
	for _, method := range tc.XMLData.Methods.Method {
		if method.Name == methodName {
			methodToExecute = &method
			break
		}
	}
	
	if methodToExecute == nil {
		return fmt.Errorf("method %s not found", methodName)
	}
	
	// Execute each step in the method
	var ifConditionResult bool
	var inIfBlock bool
	var inThenBlock bool
	var inElseBlock bool
	
	for _, step := range methodToExecute.Steps {
		if execute, ok := step.(xml.StartElement); ok {
			op := OperationType(execute.Name.Local)
			
			// Handle IF/THEN/ELSE blocks
			if inIfBlock {
				if op == OpThen {
					inThenBlock = true
					inElseBlock = false
					continue
				} else if op == OpElse {
					inThenBlock = false
					inElseBlock = true
					continue
				} else if op != OpExecute && op != OpEval && op != OpCompare && op != OpIf {
					// End of IF block
					inIfBlock = false
					inThenBlock = false
					inElseBlock = false
				}
			}
			
			// Skip operations in the wrong branch of an IF
			if inIfBlock && ((inThenBlock && !ifConditionResult) || (inElseBlock && ifConditionResult)) {
				continue
			}
			
			// Process operations
			switch op {
			case OpExecute:
				// Call another method
				methodToCall := ""
				for _, attr := range execute.Attr {
					if attr.Name.Local == "method" {
						methodToCall = attr.Value
						break
					}
				}
				if methodToCall != "" {
					if err := tc.executeMethod(methodToCall); err != nil {
						return err
					}
				}
				
			case OpEval:
				// Evaluate an expression
				var target string
				var left, right string
				var op EvalOperator
				
				for _, attr := range execute.Attr {
					switch attr.Name.Local {
					case "target":
						target = attr.Value
					case "op":
						op = EvalOperator(attr.Value)
					case "left":
						left = attr.Value
					case "right":
						right = attr.Value
					}
				}
				
				if target != "" {
					// Evaluate the expression
					result, err := tc.evaluateExpression(left, right, op)
					if err != nil {
						return fmt.Errorf("evaluation error in method %s: %w", methodName, err)
					}
					
					// Store the result in the appropriate location (input, output, or internal)
					tc.setVariableValue(target, result)
				}
				
			case OpIf:
				// Start of an IF block
				inIfBlock = true
				inThenBlock = false
				inElseBlock = false
				
				var left, right string
				var op ComparisonOperator
				
				for _, attr := range execute.Attr {
					switch attr.Name.Local {
					case "left":
						left = attr.Value
					case "right":
						right = attr.Value
					case "op":
						op = ComparisonOperator(attr.Value)
					}
				}
				
				// Evaluate the condition
				result, err := tc.evaluateComparison(left, right, op)
				if err != nil {
					return fmt.Errorf("if condition error in method %s: %w", methodName, err)
				}
				ifConditionResult = result
				
			case OpCompare:
				// Compare two values and store the result
				var target string
				var left, right string
				var op ComparisonOperator
				
				for _, attr := range execute.Attr {
					switch attr.Name.Local {
					case "target":
						target = attr.Value
					case "left":
						left = attr.Value
					case "right":
						right = attr.Value
					case "op":
						op = ComparisonOperator(attr.Value)
					}
				}
				
				if target != "" {
					// Evaluate the comparison
					result, err := tc.evaluateComparison(left, right, op)
					if err != nil {
						return fmt.Errorf("comparison error in method %s: %w", methodName, err)
					}
					
					// Store the result
					tc.setVariableValue(target, result)
				}
				
			case OpBausteinFinish:
				// End of a building block, nothing to do here
				continue
			}
		}
	}
	
	return nil
}

func (tc *TaxCalculator) getVariableValue(name string) (interface{}, error) {
	// Check inputs first
	if val, ok := tc.InputValues[name]; ok {
		return val, nil
	}

	// Then check internal variables
	if val, ok := tc.InternalVars[name]; ok {
		return val, nil
	}

	// Then check constants
	if val, ok := tc.Constants[name]; ok {
		return val, nil
	}

	// Then check outputs
	if val, ok := tc.OutputValues[name]; ok {
		return val, nil
	}

	// If it's a number literal, try to parse it
	if _, err := strconv.Atoi(name); err == nil {
		return parseValue("int", name), nil
	}

	if strings.Contains(name, ".") {
		if _, err := strconv.ParseFloat(name, 64); err == nil {
			return parseValue("double", name), nil
		}
	}

	// If it's "true" or "false"
	if name == "true" || name == "false" {
		return parseValue("boolean", name), nil
	}

	return nil, fmt.Errorf("variable %s not found", name)
}

func (tc *TaxCalculator) setVariableValue(name string, value interface{}) {
	// Set in the appropriate variable map based on the variable name
	for _, input := range tc.XMLData.Variables.Inputs.Input {
		if input.Name == name {
			tc.InputValues[name] = value
			return
		}
	}

	for _, output := range tc.XMLData.Variables.Outputs.Output {
		if output.Name == name {
			tc.OutputValues[name] = value
			return
		}
	}

	// Default to internal variables
	tc.InternalVars[name] = value
}

func (tc *TaxCalculator) evaluateComparison(left, right string, op ComparisonOperator) (bool, error) {
	leftVal, err := tc.getVariableValue(left)
	if err != nil {
		return false, fmt.Errorf("left operand error: %w", err)
	}

	rightVal, err := tc.getVariableValue(right)
	if err != nil {
		return false, fmt.Errorf("right operand error: %w", err)
	}

	// Convert to comparable types
	leftNum, rightNum, err := tc.convertToComparableNumbers(leftVal, rightVal)
	if err != nil {
		// Try boolean comparison
		leftBool, okLeft := leftVal.(bool)
		rightBool, okRight := rightVal.(bool)

		if okLeft && okRight {
			switch op {
			case CompEQ:
				return leftBool == rightBool, nil
			case CompNE:
				return leftBool != rightBool, nil
			default:
				return false, fmt.Errorf("invalid boolean comparison operator: %s", op)
			}
		}

		return false, fmt.Errorf("incompatible types for comparison: %T and %T", leftVal, rightVal)
	}

	// Perform the comparison
	switch op {
	case CompLT:
		return leftNum.Cmp(rightNum) < 0, nil
	case CompLE:
		return leftNum.Cmp(rightNum) <= 0, nil
	case CompGT:
		return leftNum.Cmp(rightNum) > 0, nil
	case CompGE:
		return leftNum.Cmp(rightNum) >= 0, nil
	case CompEQ:
		return leftNum.Cmp(rightNum) == 0, nil
	case CompNE:
		return leftNum.Cmp(rightNum) != 0, nil
	default:
		return false, fmt.Errorf("unknown comparison operator: %s", op)
	}
}

func (tc *TaxCalculator) evaluateExpression(left, right string, op EvalOperator) (interface{}, error) {
	// For simple assignment, just get the right value
	if op == EvalAssign {
		rightVal, err := tc.getVariableValue(right)
		if err != nil {
			return nil, fmt.Errorf("right operand error: %w", err)
		}
		return rightVal, nil
	}

	// For other operations, get both values
	leftVal, err := tc.getVariableValue(left)
	if err != nil {
		return nil, fmt.Errorf("left operand error: %w", err)
	}

	rightVal, err := tc.getVariableValue(right)
	if err != nil {
		return nil, fmt.Errorf("right operand error: %w", err)
	}

	// Convert to compatible numeric types
	leftNum, rightNum, err := tc.convertToCompatibleNumbers(leftVal, rightVal)
	if err != nil {
		return nil, err
	}

	// Perform the operation
	result := new(big.Rat)
	switch op {
	case EvalAdd:
		result.Add(leftNum, rightNum)
	case EvalSubtract:
		result.Sub(leftNum, rightNum)
	case EvalMultiply:
		result.Mul(leftNum, rightNum)
	case EvalDivide:
		if rightNum.Cmp(big.NewRat(0, 1)) == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result.Quo(leftNum, rightNum)
	case EvalIntDivide:
		// Integer division
		if rightNum.Cmp(big.NewRat(0, 1)) == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		// Get the integer quotient
		quo := new(big.Rat).Quo(leftNum, rightNum)
		f, _ := quo.Float64()
		intResult := int(f)
		return intResult, nil
	case EvalModulo:
		// Modulo operation
		if rightNum.Cmp(big.NewRat(0, 1)) == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		// Get the remainder of integer division
		leftF, _ := leftNum.Float64()
		rightF, _ := rightNum.Float64()
		intResult := int(leftF) % int(rightF)
		return intResult, nil
	default:
		return nil, fmt.Errorf("unknown operation: %s", op)
	}

	// For integer operations, return an int if both inputs were ints
	if _, okLeft := leftVal.(int); okLeft {
		if _, okRight := rightVal.(int); okRight {
			f, _ := result.Float64()
			return int(f), nil
		}
	}

	return result, nil
}

func (tc *TaxCalculator) convertToCompatibleNumbers(left, right interface{}) (*big.Rat, *big.Rat, error) {
	leftRat := new(big.Rat)
	rightRat := new(big.Rat)

	// Convert left value
	switch v := left.(type) {
	case int:
		leftRat.SetInt64(int64(v))
	case *big.Rat:
		leftRat.Set(v)
	case bool:
		if v {
			leftRat.SetInt64(1)
		} else {
			leftRat.SetInt64(0)
		}
	default:
		return nil, nil, fmt.Errorf("left operand is not a number: %T", left)
	}

	// Convert right value
	switch v := right.(type) {
	case int:
		rightRat.SetInt64(int64(v))
	case *big.Rat:
		rightRat.Set(v)
	case bool:
		if v {
			rightRat.SetInt64(1)
		} else {
			rightRat.SetInt64(0)
		}
	default:
		return nil, nil, fmt.Errorf("right operand is not a number: %T", right)
	}

	return leftRat, rightRat, nil
}

func (tc *TaxCalculator) convertToComparableNumbers(left, right interface{}) (*big.Rat, *big.Rat, error) {
	return tc.convertToCompatibleNumbers(left, right)
}

func parseValue(valueType string, value string) interface{} {
	switch valueType {
	case "int":
		if strings.Contains(value, "default") {
			return 0 // Default value handling
		}
		intVal, _ := strconv.Atoi(value)
		return intVal
	case "double", "BigDecimal":
		if strings.Contains(value, "BigDecimal") || strings.Contains(value, "ZERO") {
			return big.NewRat(0, 1) // Default to 0 for expressions
		}
		r := new(big.Rat)
		r.SetString(value)
		return r
	case "boolean":
		boolVal, _ := strconv.ParseBool(value)
		return boolVal
	default:
		return value
	}
}