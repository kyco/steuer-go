package models

import (
	"errors"
	"testing"
)

func TestTaxClass(t *testing.T) {
	tests := []struct {
		class    TaxClass
		expected int
	}{
		{TaxClass1, 1},
		{TaxClass2, 2},
		{TaxClass3, 3},
		{TaxClass4, 4},
		{TaxClass5, 5},
		{TaxClass6, 6},
	}

	for _, tc := range tests {
		if int(tc.class) != tc.expected {
			t.Errorf("Expected TaxClass%d to be %d, got %d", tc.expected, tc.expected, int(tc.class))
		}
	}
}

func TestPaymentPeriod(t *testing.T) {
	tests := []struct {
		period   PaymentPeriod
		expected int
	}{
		{Year, 1},
		{Month, 2},
		{Week, 3},
		{Day, 4},
	}

	for _, tc := range tests {
		if int(tc.period) != tc.expected {
			t.Errorf("Expected PaymentPeriod %d to be %d, got %d", tc.expected, tc.expected, int(tc.period))
		}
	}
}

func TestTaxRequest(t *testing.T) {
	req := TaxRequest{
		Period:   Year,
		Income:   50000,
		TaxClass: TaxClass3,
	}

	if req.Period != Year {
		t.Errorf("Expected Period %v, got %v", Year, req.Period)
	}

	if req.Income != 50000 {
		t.Errorf("Expected Income %d, got %d", 50000, req.Income)
	}

	if req.TaxClass != TaxClass3 {
		t.Errorf("Expected TaxClass %v, got %v", TaxClass3, req.TaxClass)
	}
}

func TestTaxResult(t *testing.T) {
	testErr := errors.New("test error")
	result := TaxResult{
		Income:        50000.0,
		IncomeTax:     10000.0,
		SolidarityTax: 500.0,
		TotalTax:      10500.0,
		NetIncome:     39500.0,
		TaxRate:       21.0,
		Error:         testErr,
	}

	if result.Income != 50000.0 {
		t.Errorf("Expected Income %f, got %f", 50000.0, result.Income)
	}

	if result.IncomeTax != 10000.0 {
		t.Errorf("Expected IncomeTax %f, got %f", 10000.0, result.IncomeTax)
	}

	if result.SolidarityTax != 500.0 {
		t.Errorf("Expected SolidarityTax %f, got %f", 500.0, result.SolidarityTax)
	}

	if result.TotalTax != 10500.0 {
		t.Errorf("Expected TotalTax %f, got %f", 10500.0, result.TotalTax)
	}

	if result.NetIncome != 39500.0 {
		t.Errorf("Expected NetIncome %f, got %f", 39500.0, result.NetIncome)
	}

	if result.TaxRate != 21.0 {
		t.Errorf("Expected TaxRate %f, got %f", 21.0, result.TaxRate)
	}

	if result.Error != testErr {
		t.Errorf("Expected Error %v, got %v", testErr, result.Error)
	}
}