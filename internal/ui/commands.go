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

type CalculationStartedMsg struct{}
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
		return CalculationStartedMsg{}
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
		var cmds []tea.Cmd

		incomeInCents := int(income * 100)

		taxRequest := models.TaxRequest{
			Period:   models.Year,
			Income:   incomeInCents,
			TaxClass: models.TaxClass(taxClass),
		}

		response, err := api.CalculateTax(taxRequest)

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

