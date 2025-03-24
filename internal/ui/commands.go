package ui

import (
	"fmt"
	
	tea "github.com/charmbracelet/bubbletea"
	
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/domain/models"
	"tax-calculator/internal/services"
)

// Message types

// Debug message type
type DebugLogMsg struct {
	Message string
}

// Calculation message types
type CalculationStartedMsg struct{}
type CalculationMsg struct {
	Result *api.TaxCalculationResponse
	Error  error
}

// Comparison message types
type ComparisonStartedMsg struct{}
type ComparisonProgressMsg struct {
	CompletedCalls int
	TotalCalls     int
}
type ComparisonMsg struct {
	Results []models.TaxResult
	Error   error
}

// Commands

// PerformCalculationCmd initiates the tax calculation
func PerformCalculationCmd(taxClass int, income float64, year string) tea.Cmd {
	return func() tea.Msg {
		return CalculationStartedMsg{}
	}
}

// CaptureDebugCmd creates a command to capture debug output
func CaptureDebugCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg{
			Message: message,
		}
	}
}

// FetchResultsCmd fetches the tax calculation results
func FetchResultsCmd(taxClass int, income float64) tea.Cmd {
	return func() tea.Msg {
		var cmds []tea.Cmd
		var msgs []tea.Msg
		
		// Convert income to cents
		incomeInCents := int(income * 100)

		// Create tax request
		taxRequest := models.TaxRequest{
			Period:   models.Year,
			Income:   incomeInCents,
			TaxClass: models.TaxClass(taxClass),
		}
		
		// Add more descriptive debug message before API call
		debugMsg1 := fmt.Sprintf("DEBUG: [Income: €%d.00] API request for tax class %d, period %d", 
			incomeInCents/100, taxClass, models.Year)
		msgs = append(msgs, DebugLogMsg{Message: debugMsg1})
		
		// Generate a truncated URL to show in debug log
		shortUrl := fmt.Sprintf("%s?..RE4=%d&STKL=%d", 
			api.BaseURL, incomeInCents, taxClass)
		debugMsg2 := fmt.Sprintf("DEBUG: [Income: €%d.00] URL: %s", incomeInCents/100, shortUrl)
		msgs = append(msgs, DebugLogMsg{Message: debugMsg2})

		// Call the BMF API
		response, err := api.CalculateTax(taxRequest)
		
		// Add detailed debug message after API call
		if err == nil {
			// Extract tax values for debug info
			var incomeTax, solidarityTax string
			for _, output := range response.Outputs.Output {
				if output.Name == "LSTLZZ" {
					incomeTax = output.Value
				} else if output.Name == "SOLZLZZ" {
					solidarityTax = output.Value
				}
			}
			
			// Calculate tax amounts in euros for debug info
			incomeTaxEuros := float64(api.MustParseInt(incomeTax)) / 100
			solidarityTaxEuros := float64(api.MustParseInt(solidarityTax)) / 100
			totalTax := incomeTaxEuros + solidarityTaxEuros
			
			// Print a success message with tax details
			debugMsg3 := fmt.Sprintf("DEBUG: [Income: €%d.00] ✓ Success - Tax: €%.2f (%.2f%%)", 
				incomeInCents/100, 
				totalTax,
				(totalTax / float64(incomeInCents/100)) * 100)
			msgs = append(msgs, DebugLogMsg{Message: debugMsg3})
		} else {
			// Print a failure message with error details
			debugMsg3 := fmt.Sprintf("DEBUG: [Income: €%d.00] ✗ Failed: %v", 
				incomeInCents/100, err)
			msgs = append(msgs, DebugLogMsg{Message: debugMsg3})
		}
		
		// Send all debug messages first
		for _, msg := range msgs {
			cmds = append(cmds, func() tea.Msg { return msg })
		}
		
		// Return the calculation result
		calcMsg := CalculationMsg{
			Result: response,
			Error:  err,
		}
		
		cmds = append(cmds, func() tea.Msg { return calcMsg })
		
		// Return batch of commands
		return tea.Batch(cmds...)()
	}
}

// PerformComparisonCmd initiates the comparison calculation
func PerformComparisonCmd() tea.Cmd {
	return func() tea.Msg {
		return ComparisonStartedMsg{}
	}
}

// ProgressUpdateCmd returns a command to update progress
func ProgressUpdateCmd(completed, total int) tea.Cmd {
	return func() tea.Msg {
		return ComparisonProgressMsg{
			CompletedCalls: completed,
			TotalCalls:     total,
		}
	}
}

// CompletedResultsCmd returns a command to deliver final results
func CompletedResultsCmd(results []models.TaxResult) tea.Cmd {
	return func() tea.Msg {
		return ComparisonMsg{
			Results: results,
		}
	}
}

// FetchComparisonCmd calculates taxes for multiple income points
func FetchComparisonCmd(taxClass int, income float64) tea.Cmd {
	return func() tea.Msg {
		// Create tax service
		taxService := services.NewTaxService()
		
		// Calculate half and double the income
		halfIncome := income / 2
		doubleIncome := income * 2
		
		// Initialize empty results array
		var results []models.TaxResult
		
		// Calculate tax for original income
		originalResult := calculateTaxForIncome(taxClass, income, taxService)
		results = append(results, originalResult)
		
		// Calculate tax for half income
		halfResult := calculateTaxForIncome(taxClass, halfIncome, taxService)
		results = append(results, halfResult)
		
		// Calculate tax for double income
		doubleResult := calculateTaxForIncome(taxClass, doubleIncome, taxService)
		results = append(results, doubleResult)
		
		// Calculate 10 points between half and original income
		increment := (income - halfIncome) / 10
		for i := 1; i <= 9; i++ {
			point := halfIncome + (float64(i) * increment)
			result := calculateTaxForIncome(taxClass, point, taxService)
			results = append(results, result)
		}
		
		// Calculate 10 points between original and double income
		increment = (doubleIncome - income) / 10
		for i := 1; i <= 9; i++ {
			point := income + (float64(i) * increment)
			result := calculateTaxForIncome(taxClass, point, taxService)
			results = append(results, result)
		}
		
		// Sort results by income
		sortResults(results)
		
		// Return completed results
		return ComparisonMsg{
			Results: results,
		}
	}
}

// Helper function to calculate tax for a given income point
func calculateTaxForIncome(taxClass int, income float64, taxService *services.TaxService) models.TaxResult {
	// Convert income to cents
	incomeInCents := int(income * 100)
	
	// Create tax request
	taxRequest := models.TaxRequest{
		Period:   models.Year,
		Income:   incomeInCents,
		TaxClass: models.TaxClass(taxClass),
	}
	
	// Call the BMF API
	response, err := api.CalculateTax(taxRequest)
	
	// Create result
	var result models.TaxResult
	if err != nil {
		result = models.TaxResult{
			Income: income,
			Error:  err,
		}
	} else {
		result = taxService.GetTaxSummary(response, income)
	}
	
	return result
}

// Helper function to sort results by income
func sortResults(results []models.TaxResult) {
	// Simple bubble sort
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Income > results[j].Income {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}