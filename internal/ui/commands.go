package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/domain/models"
	"tax-calculator/internal/services"
)

type DebugLogMsg struct {
	Message string
}

type CalculationStartedMsg struct{
	UseLocalCalculator bool
}
type CalculationMsg struct {
	Result *api.TaxCalculationResponse
	Error  error
}

type ComparisonStartedMsg struct{}
type ComparisonProgressMsg struct {
	CompletedCalls int
	TotalCalls     int
}
type ComparisonMsg struct {
	Results []models.TaxResult
	Error   error
}

func PerformCalculationCmd(taxClass int, income float64, year string) tea.Cmd {
	return func() tea.Msg {
		return CalculationStartedMsg{
			UseLocalCalculator: false,
		}
	}
}

func PerformCalculationWithAdvancedOptionsCmd(taxClass int, income float64, year string, advancedParams models.TaxRequest, useLocalCalculator bool) tea.Cmd {
	return func() tea.Msg {
		// Create a basic request with main parameters
		taxRequest := models.TaxRequest{
			Period:   models.Year,
			Income:   int(income * 100),
			TaxClass: models.TaxClass(taxClass),
			
			// Copy all advanced parameters
			AJAHR:  advancedParams.AJAHR,
			ALTER1: advancedParams.ALTER1,
			KRV:    advancedParams.KRV,
			KVZ:    advancedParams.KVZ,
			PVS:    advancedParams.PVS,
			PVZ:    advancedParams.PVZ,
			R:      advancedParams.R,
			ZKF:    advancedParams.ZKF,
			VBEZ:   advancedParams.VBEZ,
			VJAHR:  advancedParams.VJAHR,
			PKPV:   advancedParams.PKPV,
			PKV:    advancedParams.PKV,
			PVA:    advancedParams.PVA,
		}
		
		// Return the calculation started message first
		return tea.Batch(
			func() tea.Msg {
				return CalculationStartedMsg{
					UseLocalCalculator: useLocalCalculator,
				}
			},
			func() tea.Msg {
				// Then perform the actual calculation with advanced parameters
				return FetchResultsWithAdvancedParamsCmd(taxRequest, useLocalCalculator)()
			},
		)()
	}
}

func CaptureDebugCmd(message string) tea.Cmd {
	return func() tea.Msg {
		return DebugLogMsg{
			Message: message,
		}
	}
}

func FetchResultsCmd(taxClass int, income float64) tea.Cmd {
	return func() tea.Msg {
		return FetchResultsWithModeCmd(taxClass, income, false)()
	}
}

func FetchResultsWithModeCmd(taxClass int, income float64, useLocalCalculator bool) tea.Cmd {
	// Create a default TaxRequest with basic parameters
	taxRequest := models.TaxRequest{
		Period:   models.Year,
		Income:   int(income * 100),
		TaxClass: models.TaxClass(taxClass),
	}
	
	return FetchResultsWithAdvancedParamsCmd(taxRequest, useLocalCalculator)
}

func FetchResultsWithAdvancedParamsCmd(taxRequest models.TaxRequest, useLocalCalculator bool) tea.Cmd {
	return func() tea.Msg {
		var cmds []tea.Cmd

		taxService := services.NewTaxService()
		if useLocalCalculator {
			taxService.EnableLocalCalculator()
			cmds = append(cmds, CaptureDebugCmd("Using local tax calculator"))
		}

		// Calculate tax using the service (with local or remote calculator based on flag)
		_, err := taxService.CalculateTax(taxRequest)
		
		var response *api.TaxCalculationResponse
		if err == nil {
			if useLocalCalculator {
				// Use the local calculator directly to get the raw response
				localCalc := services.GetLocalTaxCalculator()
				if !localCalc.IsInitialized() {
					if initErr := localCalc.Initialize(); initErr != nil {
						calcMsg := CalculationMsg{
							Result: nil,
							Error:  initErr,
						}
						cmds = append(cmds, func() tea.Msg { return calcMsg })
						return tea.Batch(cmds...)()
					}
				}
				// Get the response from local calculator
				response, err = localCalc.CalculateTax(taxRequest)
			} else {
				// Use the API for remote calculation
				response, _ = api.CalculateTax(taxRequest)
			}
		}

		calcMsg := CalculationMsg{
			Result: response,
			Error:  err,
		}

		cmds = append(cmds, func() tea.Msg { return calcMsg })

		return tea.Batch(cmds...)()
	}
}

func PerformComparisonCmd() tea.Cmd {
	return func() tea.Msg {
		return ComparisonStartedMsg{}
	}
}

func ProgressUpdateCmd(completed, total int) tea.Cmd {
	return func() tea.Msg {
		return ComparisonProgressMsg{
			CompletedCalls: completed,
			TotalCalls:     total,
		}
	}
}

func CompletedResultsCmd(results []models.TaxResult) tea.Cmd {
	return func() tea.Msg {
		return ComparisonMsg{
			Results: results,
		}
	}
}

func FetchComparisonCmd(taxClass int, income float64) tea.Cmd {
	return func() tea.Msg {
		taxService := services.NewTaxService()

		halfIncome := income / 2
		doubleIncome := income * 2

		var results []models.TaxResult

		originalResult := calculateTaxForIncome(taxClass, income, taxService)
		results = append(results, originalResult)

		halfResult := calculateTaxForIncome(taxClass, halfIncome, taxService)
		results = append(results, halfResult)

		doubleResult := calculateTaxForIncome(taxClass, doubleIncome, taxService)
		results = append(results, doubleResult)

		increment := (income - halfIncome) / 10
		for i := 1; i <= 9; i++ {
			point := halfIncome + (float64(i) * increment)
			result := calculateTaxForIncome(taxClass, point, taxService)
			results = append(results, result)
		}

		increment = (doubleIncome - income) / 10
		for i := 1; i <= 9; i++ {
			point := income + (float64(i) * increment)
			result := calculateTaxForIncome(taxClass, point, taxService)
			results = append(results, result)
		}

		sortResults(results)

		return ComparisonMsg{
			Results: results,
		}
	}
}

func calculateTaxForIncome(taxClass int, income float64, taxService *services.TaxService) models.TaxResult {
	incomeInCents := int(income * 100)

	taxRequest := models.TaxRequest{
		Period:   models.Year,
		Income:   incomeInCents,
		TaxClass: models.TaxClass(taxClass),
	}

	response, err := api.CalculateTax(taxRequest)

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

func sortResults(results []models.TaxResult) {
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Income > results[j].Income {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}