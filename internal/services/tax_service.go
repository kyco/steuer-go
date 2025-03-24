package services

import (
	"fmt"
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/domain/models"
)

// TaxService provides methods for tax calculations
type TaxService struct{}

// NewTaxService creates a new TaxService
func NewTaxService() *TaxService {
	return &TaxService{}
}

// CalculateTax performs a tax calculation for a given request
func (s *TaxService) CalculateTax(req models.TaxRequest) (models.TaxResult, error) {
	response, err := api.CalculateTax(req)
	if err != nil {
		return models.TaxResult{
			Income: float64(req.Income) / 100,
			Error:  err,
		}, err
	}

	return s.GetTaxSummary(response, float64(req.Income)/100), nil
}

// GetTaxSummary extracts the key tax values from a calculation response
func (s *TaxService) GetTaxSummary(response *api.TaxCalculationResponse, income float64) models.TaxResult {
	// Initialize result with income
	result := models.TaxResult{
		Income: income,
	}

	if response == nil {
		result.Error = fmt.Errorf("no response data")
		return result
	}

	// Extract tax values
	var incomeTax, solidarityTax string
	for _, output := range response.Outputs.Output {
		if output.Name == "LSTLZZ" {
			incomeTax = output.Value
		} else if output.Name == "SOLZLZZ" {
			solidarityTax = output.Value
		}
	}

	// Convert to euros
	result.IncomeTax = float64(api.MustParseInt(incomeTax)) / 100
	result.SolidarityTax = float64(api.MustParseInt(solidarityTax)) / 100
	result.TotalTax = result.IncomeTax + result.SolidarityTax
	result.NetIncome = income - result.TotalTax
	
	// Calculate tax rate as percentage
	if income > 0 {
		result.TaxRate = (result.TotalTax / income) * 100
	}

	return result
}

// GetReadableTaxSummary returns a human-readable summary of the tax calculation
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

// CalculateComparisonTaxes calculates tax for multiple income points
func (s *TaxService) CalculateComparisonTaxes(taxClass models.TaxClass, baseIncome float64) []models.TaxResult {
	var results []models.TaxResult
	
	// Calculate half and double the income
	halfIncome := baseIncome / 2
	doubleIncome := baseIncome * 2
	
	// Create a slice to store all income points in order
	incomePoints := []float64{}
	
	// Add half income point
	incomePoints = append(incomePoints, halfIncome)
	
	// Add 10 points between half income and base income
	lowerIncrement := (baseIncome - halfIncome) / 11
	for i := 1; i <= 10; i++ {
		incomePoints = append(incomePoints, halfIncome + (lowerIncrement * float64(i)))
	}
	
	// Add base income point
	incomePoints = append(incomePoints, baseIncome)
	
	// Add 10 points between base income and double income
	higherIncrement := (doubleIncome - baseIncome) / 11
	for i := 1; i <= 10; i++ {
		incomePoints = append(incomePoints, baseIncome + (higherIncrement * float64(i)))
	}
	
	// Add double income point
	incomePoints = append(incomePoints, doubleIncome)
	
	// Process each income point synchronously
	for _, income := range incomePoints {
		// Convert income to cents
		incomeInCents := int(income * 100)
		
		// Create tax request
		taxRequest := models.TaxRequest{
			Period:   models.Year,
			Income:   incomeInCents, 
			TaxClass: taxClass,
		}
		
		// Calculate tax - API calls are made synchronously within CalculateTax
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