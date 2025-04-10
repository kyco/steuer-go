package services

import (
	"fmt"
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/domain/models"
)

type TaxService struct{
	useLocalCalculator bool
}

func NewTaxService() *TaxService {
	return &TaxService{
		useLocalCalculator: false,
	}
}

func (s *TaxService) EnableLocalCalculator() {
	s.useLocalCalculator = true
}

func (s *TaxService) DisableLocalCalculator() {
	s.useLocalCalculator = false
}

func (s *TaxService) CalculateTax(req models.TaxRequest) (models.TaxResult, error) {
	var response *api.TaxCalculationResponse
	var err error
	
	if s.useLocalCalculator {
		localCalc := GetLocalTaxCalculator()
		
		if !localCalc.IsInitialized() {
			if err := localCalc.Initialize(); err != nil {
				return models.TaxResult{
					Income: float64(req.Income) / 100,
					Error:  fmt.Errorf("failed to initialize local calculator: %w", err),
				}, err
			}
		}
		
		response, err = localCalc.CalculateTax(req)
	} else {
		response, err = api.CalculateTax(req)
		
		if err != nil {
			localCalc := GetLocalTaxCalculator()
			
			if !localCalc.IsInitialized() {
				if initErr := localCalc.Initialize(); initErr != nil {
					return models.TaxResult{
						Income: float64(req.Income) / 100,
						Error:  fmt.Errorf("API error: %v, local calculator error: %w", err, initErr),
					}, err
				}
			}
			
			response, err = localCalc.CalculateTax(req)
		}
	}

	if err != nil {
		return models.TaxResult{
			Income: float64(req.Income) / 100,
			Error:  err,
		}, err
	}

	return s.GetTaxSummary(response, float64(req.Income)/100), nil
}

func (s *TaxService) GetTaxSummary(response *api.TaxCalculationResponse, income float64) models.TaxResult {
	result := models.TaxResult{
		Income: income,
	}

	if response == nil {
		result.Error = fmt.Errorf("no response data")
		return result
	}
	var incomeTax, solidarityTax string
	for _, output := range response.Outputs.Output {
		if output.Name == "LSTLZZ" {
			incomeTax = output.Value
		} else if output.Name == "SOLZLZZ" {
			solidarityTax = output.Value
		}
	}

	result.IncomeTax = float64(api.MustParseInt(incomeTax)) / 100
	result.SolidarityTax = float64(api.MustParseInt(solidarityTax)) / 100
	result.TotalTax = result.IncomeTax + result.SolidarityTax
	result.NetIncome = income - result.TotalTax
	if income > 0 {
		result.TaxRate = (result.TotalTax / income) * 100
	}

	return result
}

func (s *TaxService) GetReadableTaxSummary(response *api.TaxCalculationResponse) string {
	var incomeTax, solidarityTax string

	for _, output := range response.Outputs.Output {
		if output.Name == "LSTLZZ" {
			incomeTax = output.Value
		} else if output.Name == "SOLZLZZ" {
			solidarityTax = output.Value
		}
	}

	incomeTaxEuros := float64(api.MustParseInt(incomeTax)) / 100
	solidarityTaxEuros := float64(api.MustParseInt(solidarityTax)) / 100
	totalTax := incomeTaxEuros + solidarityTaxEuros

	return fmt.Sprintf("Tax Summary for %s:\n"+
		"Income Tax: %.2f EUR\n"+
		"Solidarity Tax: %.2f EUR\n"+
		"Total Tax: %.2f EUR",
		response.Year, incomeTaxEuros, solidarityTaxEuros, totalTax)
}

func (s *TaxService) CalculateComparisonTaxes(taxClass models.TaxClass, baseIncome float64) []models.TaxResult {
	var results []models.TaxResult
	
	halfIncome := baseIncome / 2
	doubleIncome := baseIncome * 2
	
	incomePoints := []float64{}
	
	incomePoints = append(incomePoints, halfIncome)
	lowerIncrement := (baseIncome - halfIncome) / 11
	for i := 1; i <= 10; i++ {
		incomePoints = append(incomePoints, halfIncome + (lowerIncrement * float64(i)))
	}
	
	incomePoints = append(incomePoints, baseIncome)
	higherIncrement := (doubleIncome - baseIncome) / 11
	for i := 1; i <= 10; i++ {
		incomePoints = append(incomePoints, baseIncome + (higherIncrement * float64(i)))
	}
	
	incomePoints = append(incomePoints, doubleIncome)
	for _, income := range incomePoints {
		incomeInCents := int(income * 100)
		taxRequest := models.TaxRequest{
			Period:   models.Year,
			Income:   incomeInCents, 
			TaxClass: taxClass,
		}
		
		result, err := s.CalculateTax(taxRequest)
		if err != nil {
			result = models.TaxResult{
				Income: income,
				Error:  err,
			}
		}
		
		results = append(results, result)
	}
	
	return results
}