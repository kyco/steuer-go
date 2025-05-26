package bmf

import (
	"math/big"
	"testing"
)

func TestParseValue(t *testing.T) {
	tests := []struct {
		name      string
		valueType string
		value     string
		expected  interface{}
	}{
		{
			name:      "parse int",
			valueType: "int",
			value:     "123",
			expected:  123,
		},
		{
			name:      "parse BigDecimal",
			valueType: "BigDecimal",
			value:     "123.45",
			expected:  big.NewRat(12345, 100),
		},
		{
			name:      "parse boolean true",
			valueType: "boolean",
			value:     "true",
			expected:  true,
		},
		{
			name:      "parse boolean false",
			valueType: "boolean",
			value:     "false",
			expected:  false,
		},
		{
			name:      "parse string",
			valueType: "string",
			value:     "test",
			expected:  "test",
		},
		{
			name:      "parse unknown type as string",
			valueType: "unknown",
			value:     "test",
			expected:  "test",
		},
		{
			name:      "parse invalid int",
			valueType: "int",
			value:     "invalid",
			expected:  0,
		},
		{
			name:      "parse invalid BigDecimal",
			valueType: "BigDecimal",
			value:     "invalid",
			expected:  big.NewRat(0, 1),
		},
		{
			name:      "parse invalid boolean",
			valueType: "boolean",
			value:     "invalid",
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseValue(tt.valueType, tt.value)

			switch expected := tt.expected.(type) {
			case int:
				if result != expected {
					t.Errorf("Expected %d, got %v", expected, result)
				}
			case bool:
				if result != expected {
					t.Errorf("Expected %v, got %v", expected, result)
				}
			case string:
				if result != expected {
					t.Errorf("Expected %q, got %v", expected, result)
				}
			case *big.Rat:
				if rat, ok := result.(*big.Rat); ok {
					if rat.Cmp(expected) != 0 {
						t.Errorf("Expected %v, got %v", expected, rat)
					}
				} else {
					t.Errorf("Expected *big.Rat, got %T", result)
				}
			}
		})
	}
}

func TestNewTaxCalculator(t *testing.T) {
	papData := &PAPData{
		Variables: PAPVariables{
			Inputs: PAPInputs{
				Input: []InputVariable{
					{Name: "testInput", Type: "int", Default: "100"},
				},
			},
			Outputs: PAPOutputs{
				Output: []OutputVariable{
					{Name: "testOutput", Type: "int", Default: "0"},
				},
			},
			Internals: PAPInternals{
				Internal: []InternalVariable{
					{Name: "testInternal", Type: "int", Default: "50"},
				},
			},
		},
		Constants: PAPConstants{
			Constant: []PAPConstant{
				{Name: "testConstant", Type: "int", Value: "200"},
			},
		},
	}

	calculator := NewTaxCalculator(papData)

	if calculator.XMLData != papData {
		t.Error("XMLData not set correctly")
	}

	if calculator.InputValues["testInput"] != 100 {
		t.Errorf("Expected testInput to be 100, got %v", calculator.InputValues["testInput"])
	}

	if calculator.OutputValues["testOutput"] != 0 {
		t.Errorf("Expected testOutput to be 0, got %v", calculator.OutputValues["testOutput"])
	}

	if calculator.InternalVars["testInternal"] != 50 {
		t.Errorf("Expected testInternal to be 50, got %v", calculator.InternalVars["testInternal"])
	}

	if calculator.Constants["testConstant"] != 200 {
		t.Errorf("Expected testConstant to be 200, got %v", calculator.Constants["testConstant"])
	}
}

func TestTaxCalculatorSetInputValue(t *testing.T) {
	calculator := &TaxCalculator{
		InputValues: make(map[string]interface{}),
	}

	calculator.SetInputValue("test", 123)

	if calculator.InputValues["test"] != 123 {
		t.Errorf("Expected test to be 123, got %v", calculator.InputValues["test"])
	}
}

func TestTaxCalculatorGetOutputValue(t *testing.T) {
	calculator := &TaxCalculator{
		OutputValues: map[string]interface{}{
			"test": 456,
		},
	}

	result := calculator.GetOutputValue("test")
	if result != 456 {
		t.Errorf("Expected 456, got %v", result)
	}

	result = calculator.GetOutputValue("nonexistent")
	if result != nil {
		t.Errorf("Expected nil for nonexistent key, got %v", result)
	}
}

func TestTaxCalculatorGetVariableValue(t *testing.T) {
	calculator := &TaxCalculator{
		InputValues:  map[string]interface{}{"input1": 100},
		OutputValues: map[string]interface{}{"output1": 200},
		InternalVars: map[string]interface{}{"internal1": 300},
		Constants:    map[string]interface{}{"const1": 400},
	}

	tests := []struct {
		name     string
		variable string
		expected interface{}
		hasError bool
	}{
		{
			name:     "get input value",
			variable: "input1",
			expected: 100,
			hasError: false,
		},
		{
			name:     "get output value",
			variable: "output1",
			expected: 200,
			hasError: false,
		},
		{
			name:     "get internal value",
			variable: "internal1",
			expected: 300,
			hasError: false,
		},
		{
			name:     "get constant value",
			variable: "const1",
			expected: 400,
			hasError: false,
		},
		{
			name:     "get integer literal",
			variable: "123",
			expected: 123,
			hasError: false,
		},
		{
			name:     "get float literal",
			variable: "123.45",
			expected: big.NewRat(12345, 100),
			hasError: false,
		},
		{
			name:     "get nonexistent variable",
			variable: "nonexistent",
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.getVariableValue(tt.variable)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				switch expected := tt.expected.(type) {
				case int:
					if result != expected {
						t.Errorf("Expected %d, got %v", expected, result)
					}
				case *big.Rat:
					if rat, ok := result.(*big.Rat); ok {
						if rat.Cmp(expected) != 0 {
							t.Errorf("Expected %v, got %v", expected, rat)
						}
					} else {
						t.Errorf("Expected *big.Rat, got %T", result)
					}
				}
			}
		})
	}
}

func TestTaxCalculatorSetVariableValue(t *testing.T) {
	papData := &PAPData{
		Variables: PAPVariables{
			Inputs: PAPInputs{
				Input: []InputVariable{
					{Name: "input1", Type: "int"},
				},
			},
			Outputs: PAPOutputs{
				Output: []OutputVariable{
					{Name: "output1", Type: "int"},
				},
			},
			Internals: PAPInternals{
				Internal: []InternalVariable{
					{Name: "internal1", Type: "int"},
				},
			},
		},
	}

	calculator := NewTaxCalculator(papData)

	calculator.setVariableValue("input1", 100)
	if calculator.InputValues["input1"] != 100 {
		t.Errorf("Expected input1 to be 100, got %v", calculator.InputValues["input1"])
	}

	calculator.setVariableValue("output1", 200)
	if calculator.OutputValues["output1"] != 200 {
		t.Errorf("Expected output1 to be 200, got %v", calculator.OutputValues["output1"])
	}

	calculator.setVariableValue("internal1", 300)
	if calculator.InternalVars["internal1"] != 300 {
		t.Errorf("Expected internal1 to be 300, got %v", calculator.InternalVars["internal1"])
	}

	calculator.setVariableValue("unknown", 400)
	if calculator.InternalVars["unknown"] != 400 {
		t.Errorf("Expected unknown to be stored as internal variable with value 400, got %v", calculator.InternalVars["unknown"])
	}
}

func TestTaxCalculatorConvertToCompatibleNumbers(t *testing.T) {
	calculator := &TaxCalculator{}

	tests := []struct {
		name      string
		left      interface{}
		right     interface{}
		expectErr bool
	}{
		{
			name:      "int and int",
			left:      100,
			right:     200,
			expectErr: false,
		},
		{
			name:      "int and big.Rat",
			left:      100,
			right:     big.NewRat(200, 1),
			expectErr: false,
		},
		{
			name:      "big.Rat and int",
			left:      big.NewRat(100, 1),
			right:     200,
			expectErr: false,
		},
		{
			name:      "big.Rat and big.Rat",
			left:      big.NewRat(100, 1),
			right:     big.NewRat(200, 1),
			expectErr: false,
		},
		{
			name:      "bool and int",
			left:      true,
			right:     100,
			expectErr: false,
		},
		{
			name:      "string and int",
			left:      "invalid",
			right:     100,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftRat, rightRat, err := calculator.convertToCompatibleNumbers(tt.left, tt.right)

			if tt.expectErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if leftRat == nil || rightRat == nil {
					t.Error("Expected non-nil rational numbers")
				}
			}
		})
	}
}

func TestTaxCalculatorEvaluateComparison(t *testing.T) {
	calculator := &TaxCalculator{
		InputValues: map[string]interface{}{
			"val1": 100,
			"val2": 200,
		},
	}

	tests := []struct {
		name     string
		left     string
		right    string
		op       ComparisonOperator
		expected bool
		hasError bool
	}{
		{
			name:     "less than true",
			left:     "val1",
			right:    "val2",
			op:       CompLT,
			expected: true,
			hasError: false,
		},
		{
			name:     "less than false",
			left:     "val2",
			right:    "val1",
			op:       CompLT,
			expected: false,
			hasError: false,
		},
		{
			name:     "less than or equal true",
			left:     "val1",
			right:    "val1",
			op:       CompLE,
			expected: true,
			hasError: false,
		},
		{
			name:     "greater than true",
			left:     "val2",
			right:    "val1",
			op:       CompGT,
			expected: true,
			hasError: false,
		},
		{
			name:     "greater than or equal true",
			left:     "val2",
			right:    "val2",
			op:       CompGE,
			expected: true,
			hasError: false,
		},
		{
			name:     "equal true",
			left:     "val1",
			right:    "100",
			op:       CompEQ,
			expected: true,
			hasError: false,
		},
		{
			name:     "not equal true",
			left:     "val1",
			right:    "val2",
			op:       CompNE,
			expected: true,
			hasError: false,
		},
		{
			name:     "invalid left operand",
			left:     "nonexistent",
			right:    "val1",
			op:       CompEQ,
			expected: false,
			hasError: true,
		},
		{
			name:     "invalid right operand",
			left:     "val1",
			right:    "nonexistent",
			op:       CompEQ,
			expected: false,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.evaluateComparison(tt.left, tt.right, tt.op)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				if result != tt.expected {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestTaxCalculatorEvaluateExpression(t *testing.T) {
	calculator := &TaxCalculator{
		InputValues: map[string]interface{}{
			"val1": 100,
			"val2": 50,
		},
	}

	tests := []struct {
		name     string
		left     string
		right    string
		op       EvalOperator
		expected interface{}
		hasError bool
	}{
		{
			name:     "assignment",
			left:     "val1",
			right:    "val1",
			op:       EvalAssign,
			expected: 100,
			hasError: false,
		},
		{
			name:     "addition",
			left:     "val1",
			right:    "val2",
			op:       EvalAdd,
			expected: 150,
			hasError: false,
		},
		{
			name:     "subtraction",
			left:     "val1",
			right:    "val2",
			op:       EvalSubtract,
			expected: 50,
			hasError: false,
		},
		{
			name:     "multiplication",
			left:     "val1",
			right:    "val2",
			op:       EvalMultiply,
			expected: 5000,
			hasError: false,
		},
		{
			name:     "division",
			left:     "val1",
			right:    "val2",
			op:       EvalDivide,
			expected: 2,
			hasError: false,
		},
		{
			name:     "integer division",
			left:     "val1",
			right:    "val2",
			op:       EvalIntDivide,
			expected: 2,
			hasError: false,
		},
		{
			name:     "modulo",
			left:     "val1",
			right:    "val2",
			op:       EvalModulo,
			expected: 0,
			hasError: false,
		},
		{
			name:     "division by zero",
			left:     "val1",
			right:    "0",
			op:       EvalDivide,
			expected: nil,
			hasError: true,
		},
		{
			name:     "invalid left operand",
			left:     "nonexistent",
			right:    "val1",
			op:       EvalAdd,
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := calculator.evaluateExpression(tt.left, tt.right, tt.op)

			if tt.hasError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}

				switch expected := tt.expected.(type) {
				case int:
					if result != expected {
						t.Errorf("Expected %d, got %v", expected, result)
					}
				case *big.Rat:
					if rat, ok := result.(*big.Rat); ok {
						if rat.Cmp(expected) != 0 {
							t.Errorf("Expected %v, got %v", expected, rat)
						}
					} else {
						t.Errorf("Expected *big.Rat, got %T", result)
					}
				}
			}
		})
	}
}

func TestOperationTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant OperationType
		expected string
	}{
		{"OpExecute", OpExecute, "EXECUTE"},
		{"OpEval", OpEval, "EVAL"},
		{"OpIf", OpIf, "IF"},
		{"OpCompare", OpCompare, "COMPARE"},
		{"OpBausteinFinish", OpBausteinFinish, "BAUSTEINFINISH"},
		{"OpThen", OpThen, "THEN"},
		{"OpElse", OpElse, "ELSE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(tt.constant))
			}
		})
	}
}

func TestComparisonOperatorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant ComparisonOperator
		expected string
	}{
		{"CompLT", CompLT, "LT"},
		{"CompLE", CompLE, "LE"},
		{"CompGT", CompGT, "GT"},
		{"CompGE", CompGE, "GE"},
		{"CompEQ", CompEQ, "EQ"},
		{"CompNE", CompNE, "NE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(tt.constant))
			}
		})
	}
}

func TestEvalOperatorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant EvalOperator
		expected string
	}{
		{"EvalAssign", EvalAssign, "="},
		{"EvalAdd", EvalAdd, "+"},
		{"EvalSubtract", EvalSubtract, "-"},
		{"EvalMultiply", EvalMultiply, "*"},
		{"EvalDivide", EvalDivide, "/"},
		{"EvalIntDivide", EvalIntDivide, "DIV"},
		{"EvalModulo", EvalModulo, "MOD"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(tt.constant))
			}
		})
	}
}

func TestXMLStructs(t *testing.T) {
	// Test that XML structs can be created and have expected fields
	papData := PAPData{
		Name:    "test",
		Version: "1.0",
	}

	if papData.Name != "test" {
		t.Errorf("Expected Name 'test', got %q", papData.Name)
	}

	if papData.Version != "1.0" {
		t.Errorf("Expected Version '1.0', got %q", papData.Version)
	}

	inputVar := InputVariable{
		Name:    "testInput",
		Type:    "int",
		Default: "0",
	}

	if inputVar.Name != "testInput" {
		t.Errorf("Expected Name 'testInput', got %q", inputVar.Name)
	}

	outputVar := OutputVariable{
		Name:    "testOutput",
		Type:    "int",
		Default: "0",
	}

	if outputVar.Name != "testOutput" {
		t.Errorf("Expected Name 'testOutput', got %q", outputVar.Name)
	}

	internalVar := InternalVariable{
		Name:    "testInternal",
		Type:    "int",
		Default: "0",
	}

	if internalVar.Name != "testInternal" {
		t.Errorf("Expected Name 'testInternal', got %q", internalVar.Name)
	}

	constant := PAPConstant{
		Name:  "testConstant",
		Type:  "int",
		Value: "100",
	}

	if constant.Name != "testConstant" {
		t.Errorf("Expected Name 'testConstant', got %q", constant.Name)
	}
}
