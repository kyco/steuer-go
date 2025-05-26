package components

import (
	"strings"
	"testing"
)

func TestNewResultsDashboard(t *testing.T) {
	result := TaxResult{
		Income:           50000,
		IncomeTax:        8000,
		SolidarityTax:    400,
		TotalTax:         8400,
		NetIncome:        41600,
		EffectiveTaxRate: 16.8,
		TaxClass:         1,
		Year:             "2024",
	}

	rd := NewResultsDashboard(result)

	if rd.Result.Income != result.Income {
		t.Errorf("Expected Income %.2f, got %.2f", result.Income, rd.Result.Income)
	}

	if rd.Result.IncomeTax != result.IncomeTax {
		t.Errorf("Expected IncomeTax %.2f, got %.2f", result.IncomeTax, rd.Result.IncomeTax)
	}

	if rd.Result.SolidarityTax != result.SolidarityTax {
		t.Errorf("Expected SolidarityTax %.2f, got %.2f", result.SolidarityTax, rd.Result.SolidarityTax)
	}

	if rd.Result.TotalTax != result.TotalTax {
		t.Errorf("Expected TotalTax %.2f, got %.2f", result.TotalTax, rd.Result.TotalTax)
	}

	if rd.Result.NetIncome != result.NetIncome {
		t.Errorf("Expected NetIncome %.2f, got %.2f", result.NetIncome, rd.Result.NetIncome)
	}

	if rd.Result.EffectiveTaxRate != result.EffectiveTaxRate {
		t.Errorf("Expected EffectiveTaxRate %.2f, got %.2f", result.EffectiveTaxRate, rd.Result.EffectiveTaxRate)
	}

	if rd.Result.TaxClass != result.TaxClass {
		t.Errorf("Expected TaxClass %d, got %d", result.TaxClass, rd.Result.TaxClass)
	}

	if rd.Result.Year != result.Year {
		t.Errorf("Expected Year %q, got %q", result.Year, rd.Result.Year)
	}

	if rd.ShowComparison {
		t.Error("Expected ShowComparison to be false by default")
	}

	if rd.Width != 80 {
		t.Errorf("Expected default Width 80, got %d", rd.Width)
	}

	if len(rd.ComparisonData) != 0 {
		t.Errorf("Expected empty ComparisonData, got %d items", len(rd.ComparisonData))
	}
}

func TestResultsDashboardView(t *testing.T) {
	result := TaxResult{
		Income:           50000,
		IncomeTax:        8000,
		SolidarityTax:    400,
		TotalTax:         8400,
		NetIncome:        41600,
		EffectiveTaxRate: 16.8,
		TaxClass:         1,
		Year:             "2024",
	}

	rd := NewResultsDashboard(result)
	view := rd.View()

	if view == "" {
		t.Error("View should not be empty")
	}

	if !strings.Contains(view, "Tax Calculation Results") {
		t.Error("View should contain header")
	}

	if !strings.Contains(view, "2024") {
		t.Error("View should contain year")
	}

	if !strings.Contains(view, "€ 50000.00") {
		t.Error("View should contain formatted income")
	}

	if !strings.Contains(view, "€ 8400.00") {
		t.Error("View should contain formatted total tax")
	}

	if !strings.Contains(view, "€ 41600.00") {
		t.Error("View should contain formatted net income")
	}

	if !strings.Contains(view, "16.8%") {
		t.Error("View should contain formatted tax rate")
	}
}

func TestResultsDashboardRenderHeader(t *testing.T) {
	result := TaxResult{Year: "2024"}
	rd := NewResultsDashboard(result)

	header := rd.renderHeader()
	if !strings.Contains(header, "Tax Calculation Results") {
		t.Error("Header should contain title")
	}

	if !strings.Contains(header, "2024") {
		t.Error("Header should contain year")
	}

	if !strings.Contains(header, "🧮") {
		t.Error("Header should contain calculator emoji")
	}
}

func TestResultsDashboardRenderKeyMetrics(t *testing.T) {
	result := TaxResult{
		Income:           50000,
		TotalTax:         8400,
		NetIncome:        41600,
		EffectiveTaxRate: 16.8,
	}
	rd := NewResultsDashboard(result)

	metrics := rd.renderKeyMetrics()
	if !strings.Contains(metrics, "Gross Income") {
		t.Error("Metrics should contain gross income label")
	}

	if !strings.Contains(metrics, "Total Tax") {
		t.Error("Metrics should contain total tax label")
	}

	if !strings.Contains(metrics, "Net Income") {
		t.Error("Metrics should contain net income label")
	}

	if !strings.Contains(metrics, "€ 50000.00") {
		t.Error("Metrics should contain formatted income")
	}

	if !strings.Contains(metrics, "€ 8400.00") {
		t.Error("Metrics should contain formatted tax")
	}

	if !strings.Contains(metrics, "€ 41600.00") {
		t.Error("Metrics should contain formatted net income")
	}
}

func TestResultsDashboardRenderVisualBreakdown(t *testing.T) {
	result := TaxResult{
		Income:        50000,
		IncomeTax:     8000,
		SolidarityTax: 400,
		NetIncome:     41600,
	}
	rd := NewResultsDashboard(result)

	breakdown := rd.renderVisualBreakdown()
	if !strings.Contains(breakdown, "Tax Breakdown") {
		t.Error("Breakdown should contain title")
	}

	if !strings.Contains(breakdown, "Income Tax") {
		t.Error("Breakdown should contain income tax label")
	}

	if !strings.Contains(breakdown, "Solidarity Tax") {
		t.Error("Breakdown should contain solidarity tax label")
	}

	if !strings.Contains(breakdown, "Net Income") {
		t.Error("Breakdown should contain net income label")
	}

	if !strings.Contains(breakdown, "━") {
		t.Error("Breakdown should contain visual bars")
	}
}

func TestResultsDashboardRenderBreakdownRow(t *testing.T) {
	result := TaxResult{}
	rd := NewResultsDashboard(result)

	row := rd.renderBreakdownRow("Test Label", 1000, 20.0, 50, "#FF0000")
	if !strings.Contains(row, "Test Label") {
		t.Error("Row should contain label")
	}

	if !strings.Contains(row, "€ 1000.00") {
		t.Error("Row should contain formatted amount")
	}

	if !strings.Contains(row, "20.0%") {
		t.Error("Row should contain formatted percentage")
	}

	if !strings.Contains(row, "━") {
		t.Error("Row should contain filled bar characters")
	}

	if !strings.Contains(row, "─") {
		t.Error("Row should contain empty bar characters")
	}
}

func TestResultsDashboardRenderMonthlyBreakdown(t *testing.T) {
	result := TaxResult{
		Income:    60000,
		TotalTax:  12000,
		NetIncome: 48000,
	}
	rd := NewResultsDashboard(result)

	monthly := rd.renderMonthlyBreakdown()
	if !strings.Contains(monthly, "Monthly Breakdown") {
		t.Error("Monthly breakdown should contain title")
	}

	if !strings.Contains(monthly, "Gross Monthly:") {
		t.Error("Monthly breakdown should contain gross monthly label")
	}

	if !strings.Contains(monthly, "Tax Monthly:") {
		t.Error("Monthly breakdown should contain tax monthly label")
	}

	if !strings.Contains(monthly, "Net Monthly:") {
		t.Error("Monthly breakdown should contain net monthly label")
	}

	if !strings.Contains(monthly, "Daily Net:") {
		t.Error("Monthly breakdown should contain daily net label")
	}

	if !strings.Contains(monthly, "€ 5000.00") {
		t.Error("Monthly breakdown should contain monthly gross (60000/12)")
	}

	if !strings.Contains(monthly, "€ 1000.00") {
		t.Error("Monthly breakdown should contain monthly tax (12000/12)")
	}

	if !strings.Contains(monthly, "€ 4000.00") {
		t.Error("Monthly breakdown should contain monthly net (48000/12)")
	}
}

func TestResultsDashboardRenderComparison(t *testing.T) {
	result := TaxResult{
		Income:           50000,
		EffectiveTaxRate: 16.8,
	}
	rd := NewResultsDashboard(result)

	comparisonData := []TaxResult{
		{Income: 30000, EffectiveTaxRate: 12.0},
		{Income: 50000, EffectiveTaxRate: 16.8},
		{Income: 70000, EffectiveTaxRate: 22.0},
	}
	rd.SetComparison(comparisonData)

	comparison := rd.renderComparison()
	if !strings.Contains(comparison, "Tax Rate Comparison") {
		t.Error("Comparison should contain title")
	}

	if !strings.Contains(comparison, "Income") {
		t.Error("Comparison should contain income header")
	}

	if !strings.Contains(comparison, "Tax Rate") {
		t.Error("Comparison should contain tax rate header")
	}

	if !strings.Contains(comparison, "€ 30000.00") {
		t.Error("Comparison should contain first income")
	}

	if !strings.Contains(comparison, "€ 50000.00") {
		t.Error("Comparison should contain second income")
	}

	if !strings.Contains(comparison, "€ 70000.00") {
		t.Error("Comparison should contain third income")
	}

	if !strings.Contains(comparison, "← Your Income") {
		t.Error("Comparison should highlight current income")
	}
}

func TestResultsDashboardSetComparison(t *testing.T) {
	result := TaxResult{}
	rd := NewResultsDashboard(result)

	if rd.ShowComparison {
		t.Error("Expected ShowComparison to be false initially")
	}

	comparisonData := []TaxResult{
		{Income: 30000, EffectiveTaxRate: 12.0},
		{Income: 50000, EffectiveTaxRate: 16.8},
	}

	rd.SetComparison(comparisonData)

	if !rd.ShowComparison {
		t.Error("Expected ShowComparison to be true after SetComparison")
	}

	if len(rd.ComparisonData) != 2 {
		t.Errorf("Expected 2 comparison items, got %d", len(rd.ComparisonData))
	}

	if rd.ComparisonData[0].Income != 30000 {
		t.Errorf("Expected first comparison income 30000, got %.2f", rd.ComparisonData[0].Income)
	}

	if rd.ComparisonData[1].Income != 50000 {
		t.Errorf("Expected second comparison income 50000, got %.2f", rd.ComparisonData[1].Income)
	}
}

func TestResultsDashboardToggleComparison(t *testing.T) {
	result := TaxResult{}
	rd := NewResultsDashboard(result)

	if rd.ShowComparison {
		t.Error("Expected ShowComparison to be false initially")
	}

	rd.ToggleComparison()
	if !rd.ShowComparison {
		t.Error("Expected ShowComparison to be true after first toggle")
	}

	rd.ToggleComparison()
	if rd.ShowComparison {
		t.Error("Expected ShowComparison to be false after second toggle")
	}
}

func TestResultsDashboardViewWithComparison(t *testing.T) {
	result := TaxResult{
		Income:           50000,
		IncomeTax:        8000,
		SolidarityTax:    400,
		TotalTax:         8400,
		NetIncome:        41600,
		EffectiveTaxRate: 16.8,
		TaxClass:         1,
		Year:             "2024",
	}

	rd := NewResultsDashboard(result)
	comparisonData := []TaxResult{
		{Income: 30000, EffectiveTaxRate: 12.0},
		{Income: 50000, EffectiveTaxRate: 16.8},
	}
	rd.SetComparison(comparisonData)

	view := rd.View()
	if !strings.Contains(view, "Tax Rate Comparison") {
		t.Error("View should contain comparison section when comparison data is set")
	}
}

func TestFormatEuro(t *testing.T) {
	tests := []struct {
		amount   float64
		expected string
	}{
		{1000.00, "€ 1000.00"},
		{1234.56, "€ 1234.56"},
		{0.00, "€ 0.00"},
		{999.99, "€ 999.99"},
		{50000.00, "€ 50000.00"},
	}

	for _, tt := range tests {
		result := formatEuro(tt.amount)
		if result != tt.expected {
			t.Errorf("formatEuro(%.2f) = %q, expected %q", tt.amount, result, tt.expected)
		}
	}
}

func TestFormatPercent(t *testing.T) {
	tests := []struct {
		percent  float64
		expected string
	}{
		{16.8, "16.8%"},
		{0.0, "0.0%"},
		{100.0, "100.0%"},
		{12.34, "12.3%"},
		{99.99, "100.0%"},
	}

	for _, tt := range tests {
		result := formatPercent(tt.percent)
		if result != tt.expected {
			t.Errorf("formatPercent(%.2f) = %q, expected %q", tt.percent, result, tt.expected)
		}
	}
}

func TestTaxResultStruct(t *testing.T) {
	result := TaxResult{
		Income:           50000,
		IncomeTax:        8000,
		SolidarityTax:    400,
		TotalTax:         8400,
		NetIncome:        41600,
		EffectiveTaxRate: 16.8,
		TaxClass:         1,
		Year:             "2024",
	}

	if result.Income != 50000 {
		t.Errorf("Expected Income 50000, got %.2f", result.Income)
	}

	if result.IncomeTax != 8000 {
		t.Errorf("Expected IncomeTax 8000, got %.2f", result.IncomeTax)
	}

	if result.SolidarityTax != 400 {
		t.Errorf("Expected SolidarityTax 400, got %.2f", result.SolidarityTax)
	}

	if result.TotalTax != 8400 {
		t.Errorf("Expected TotalTax 8400, got %.2f", result.TotalTax)
	}

	if result.NetIncome != 41600 {
		t.Errorf("Expected NetIncome 41600, got %.2f", result.NetIncome)
	}

	if result.EffectiveTaxRate != 16.8 {
		t.Errorf("Expected EffectiveTaxRate 16.8, got %.2f", result.EffectiveTaxRate)
	}

	if result.TaxClass != 1 {
		t.Errorf("Expected TaxClass 1, got %d", result.TaxClass)
	}

	if result.Year != "2024" {
		t.Errorf("Expected Year '2024', got %q", result.Year)
	}
}
