package calculation

import (
	"errors"
	"fmt"
	"tax-calculator/internal/tax/bmf"
	"tax-calculator/internal/tax/models"
	"testing"
)

// Mock API response for testing
func mockTaxResponse(incomeTax, solidarityTax string) *bmf.TaxCalculationResponse {
	return &bmf.TaxCalculationResponse{
		Year:        "2025",
		Information: "Mock tax response",
		Inputs: bmf.Inputs{
			Input: []bmf.Input{
				{Name: "LZZ", Value: "1", Status: "ok"},
				{Name: "RE4", Value: "50000", Status: "ok"},
				{Name: "STKL", Value: "1", Status: "ok"},
			},
		},
		Outputs: bmf.Outputs{
			Output: []bmf.Output{
				{Name: "LSTLZZ", Value: incomeTax, Type: "STANDARD"},
				{Name: "SOLZLZZ", Value: solidarityTax, Type: "STANDARD"},
			},
		},
	}
}

// Define interface for API to enable mocking
type TaxCalculator interface {
	CalculateTax(req models.TaxRequest) (*bmf.TaxCalculationResponse, error)
}

// Mock implementation of tax calculator
type MockTaxCalculator struct {
	ShouldFail bool
	TaxRate    float64
	SoliRate   float64
}

func (m *MockTaxCalculator) CalculateTax(req models.TaxRequest) (*bmf.TaxCalculationResponse, error) {
	if m.ShouldFail {
		return nil, errors.New("mock API error")
	}

	incomeCents := req.Income
	incomeTaxCents := int(float64(incomeCents) * m.TaxRate)
	solidarityCents := int(float64(incomeCents) * m.SoliRate)

	return mockTaxResponse(fmt.Sprintf("%d", incomeTaxCents), fmt.Sprintf("%d", solidarityCents)), nil
}

func TestNewTaxService(t *testing.T) {
	service := NewTaxService()
	if service == nil {
		t.Error("Expected non-nil TaxService")
	}
}

func TestGetTaxSummary(t *testing.T) {
	service := NewTaxService()
	tests := []struct {
		name        string
		response    *bmf.TaxCalculationResponse
		income      float64
		expected    models.TaxResult
		expectError bool
	}{
		{
			name:     "Valid response",
			response: mockTaxResponse("800000", "40000"),
			income:   50000.0,
			expected: models.TaxResult{
				Income:        50000.0,
				IncomeTax:     8000.0,
				SolidarityTax: 400.0,
				TotalTax:      8400.0,
				NetIncome:     41600.0,
				TaxRate:       16.8,
				Error:         nil,
			},
			expectError: false,
		},
		{
			name:     "Nil response",
			response: nil,
			income:   50000.0,
			expected: models.TaxResult{
				Income: 50000.0,
				Error:  fmt.Errorf("no response data"),
			},
			expectError: true,
		},
		{
			name:     "Zero income",
			response: mockTaxResponse("0", "0"),
			income:   0.0,
			expected: models.TaxResult{
				Income:        0.0,
				IncomeTax:     0.0,
				SolidarityTax: 0.0,
				TotalTax:      0.0,
				NetIncome:     0.0,
				TaxRate:       0.0,
				Error:         nil,
			},
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := service.GetTaxSummary(tc.response, tc.income)

			if tc.expectError && result.Error == nil {
				t.Errorf("Expected error, got nil")
			}

			if !tc.expectError && result.Error != nil {
				t.Errorf("Expected no error, got: %v", result.Error)
			}

			if tc.expected.Income != result.Income {
				t.Errorf("Income: expected %f, got %f", tc.expected.Income, result.Income)
			}

			if !tc.expectError {
				if tc.expected.IncomeTax != result.IncomeTax {
					t.Errorf("IncomeTax: expected %f, got %f", tc.expected.IncomeTax, result.IncomeTax)
				}

				if tc.expected.SolidarityTax != result.SolidarityTax {
					t.Errorf("SolidarityTax: expected %f, got %f", tc.expected.SolidarityTax, result.SolidarityTax)
				}

				if tc.expected.TotalTax != result.TotalTax {
					t.Errorf("TotalTax: expected %f, got %f", tc.expected.TotalTax, result.TotalTax)
				}

				if tc.expected.NetIncome != result.NetIncome {
					t.Errorf("NetIncome: expected %f, got %f", tc.expected.NetIncome, result.NetIncome)
				}

				if tc.expected.TaxRate != result.TaxRate {
					t.Errorf("TaxRate: expected %f, got %f", tc.expected.TaxRate, result.TaxRate)
				}
			}
		})
	}
}

func TestGetReadableTaxSummary(t *testing.T) {
	service := NewTaxService()
	response := mockTaxResponse("800000", "40000")

	summary := service.GetReadableTaxSummary(response)
	expected := "Tax Summary for 2025:\nIncome Tax: 8000.00 EUR\nSolidarity Tax: 400.00 EUR\nTotal Tax: 8400.00 EUR"

	if summary != expected {
		t.Errorf("Expected summary:\n%s\nGot:\n%s", expected, summary)
	}
}

// Using our own implementation to test without needing to modify the package function
func TestCalculateTaxWithMock(t *testing.T) {
	// Create a wrapper for CalculateTax that uses our mock
	calcWithMock := func(calculator TaxCalculator, req models.TaxRequest) (models.TaxResult, error) {
		response, err := calculator.CalculateTax(req)
		if err != nil {
			return models.TaxResult{
				Income: float64(req.Income) / 100,
				Error:  err,
			}, err
		}

		service := NewTaxService()
		return service.GetTaxSummary(response, float64(req.Income)/100), nil
	}

	// Test success case
	t.Run("Success", func(t *testing.T) {
		mock := &MockTaxCalculator{
			ShouldFail: false,
			TaxRate:    0.16,  // 16% tax
			SoliRate:   0.008, // 0.8% solidarity
		}

		req := models.TaxRequest{
			Period:   models.Year,
			Income:   5000000, // 50,000.00 EUR
			TaxClass: models.TaxClass1,
		}

		result, err := calcWithMock(mock, req)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		expectedIncomeTax := 8000.0 // 16% of 50000
		if result.IncomeTax != expectedIncomeTax {
			t.Errorf("Expected income tax %f, got %f", expectedIncomeTax, result.IncomeTax)
		}

		expectedSolidarityTax := 400.0 // 0.8% of 50000
		if result.SolidarityTax != expectedSolidarityTax {
			t.Errorf("Expected solidarity tax %f, got %f", expectedSolidarityTax, result.SolidarityTax)
		}
	})

	// Test error case
	t.Run("Error", func(t *testing.T) {
		mock := &MockTaxCalculator{
			ShouldFail: true,
		}

		req := models.TaxRequest{
			Period:   models.Year,
			Income:   5000000, // 50,000.00
			TaxClass: models.TaxClass1,
		}

		result, err := calcWithMock(mock, req)

		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if result.Income != 50000.0 {
			t.Errorf("Expected income 50000.0, got %f", result.Income)
		}

		if result.Error == nil {
			t.Error("Expected non-nil error in result")
		}
	})
}

func TestCalculateComparisonTaxes(t *testing.T) {
	// Since we can't modify the package function, we'll create a similar
	// function that accepts a calculator
	calculateComparisonWithMock := func(calculator TaxCalculator, taxClass models.TaxClass, baseIncome float64) []models.TaxResult {
		var results []models.TaxResult

		halfIncome := baseIncome / 2
		doubleIncome := baseIncome * 2

		incomePoints := []float64{}

		incomePoints = append(incomePoints, halfIncome)
		lowerIncrement := (baseIncome - halfIncome) / 11
		for i := 1; i <= 10; i++ {
			incomePoints = append(incomePoints, halfIncome+(lowerIncrement*float64(i)))
		}

		incomePoints = append(incomePoints, baseIncome)
		higherIncrement := (doubleIncome - baseIncome) / 11
		for i := 1; i <= 10; i++ {
			incomePoints = append(incomePoints, baseIncome+(higherIncrement*float64(i)))
		}

		incomePoints = append(incomePoints, doubleIncome)
		service := NewTaxService()

		for _, income := range incomePoints {
			incomeInCents := int(income * 100)
			taxRequest := models.TaxRequest{
				Period:   models.Year,
				Income:   incomeInCents,
				TaxClass: taxClass,
			}

			response, err := calculator.CalculateTax(taxRequest)
			var result models.TaxResult

			if err != nil {
				result = models.TaxResult{
					Income: income,
					Error:  err,
				}
			} else {
				result = service.GetTaxSummary(response, income)
			}

			results = append(results, result)
		}

		return results
	}

	mock := &MockTaxCalculator{
		ShouldFail: false,
		TaxRate:    0.20, // 20% tax
		SoliRate:   0.01, // 1% solidarity
	}

	results := calculateComparisonWithMock(mock, models.TaxClass1, 50000.0)

	// We should have 23 results (base income + 10 lower + 10 higher + half and double)
	if len(results) != 23 {
		t.Errorf("Expected 23 results, got %d", len(results))
	}

	// Check the first result (half income)
	if results[0].Income != 25000.0 {
		t.Errorf("First result income: expected 25000.0, got %f", results[0].Income)
	}

	// Check the middle result (base income)
	if results[11].Income != 50000.0 {
		t.Errorf("Middle result income: expected 50000.0, got %f", results[11].Income)
	}

	// Check the last result (double income)
	if results[22].Income != 100000.0 {
		t.Errorf("Last result income: expected 100000.0, got %f", results[22].Income)
	}

	// Check tax calculation for a sample result
	baseIncomeResult := results[11] // The one with exactly the base income
	expectedTaxRate := 21.0         // 20% income tax + 1% solidarity

	if baseIncomeResult.TaxRate != expectedTaxRate {
		t.Errorf("Tax rate: expected %f, got %f", expectedTaxRate, baseIncomeResult.TaxRate)
	}
}

// Test error handling in CalculateComparisonTaxes
func TestCalculateComparisonTaxesError(t *testing.T) {
	// Since we can't modify the package function, we'll create a similar
	// function that accepts a calculator
	calculateComparisonWithMock := func(calculator TaxCalculator, taxClass models.TaxClass, baseIncome float64) []models.TaxResult {
		var results []models.TaxResult

		halfIncome := baseIncome / 2
		doubleIncome := baseIncome * 2

		incomePoints := []float64{}

		incomePoints = append(incomePoints, halfIncome)
		lowerIncrement := (baseIncome - halfIncome) / 11
		for i := 1; i <= 10; i++ {
			incomePoints = append(incomePoints, halfIncome+(lowerIncrement*float64(i)))
		}

		incomePoints = append(incomePoints, baseIncome)
		higherIncrement := (doubleIncome - baseIncome) / 11
		for i := 1; i <= 10; i++ {
			incomePoints = append(incomePoints, baseIncome+(higherIncrement*float64(i)))
		}

		incomePoints = append(incomePoints, doubleIncome)

		for _, income := range incomePoints {
			incomeInCents := int(income * 100)
			taxRequest := models.TaxRequest{
				Period:   models.Year,
				Income:   incomeInCents,
				TaxClass: taxClass,
			}

			response, err := calculator.CalculateTax(taxRequest)
			var result models.TaxResult

			if err != nil {
				result = models.TaxResult{
					Income: income,
					Error:  err,
				}
			} else {
				service := NewTaxService()
				result = service.GetTaxSummary(response, income)
			}

			results = append(results, result)
		}

		return results
	}

	mock := &MockTaxCalculator{
		ShouldFail: true,
	}

	results := calculateComparisonWithMock(mock, models.TaxClass1, 50000.0)

	// We should still have 23 results, but all with errors
	if len(results) != 23 {
		t.Errorf("Expected 23 results, got %d", len(results))
	}

	// Check that all results have errors
	for i, result := range results {
		if result.Error == nil {
			t.Errorf("Result %d: expected error, got nil", i)
		}
	}
}
