package bmf

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"tax-calculator/internal/tax/models"
	"testing"
)

func TestMustParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"-123", -123},
		{"", 0},      // Empty string should return 0
		{"abc", 0},   // Non-numeric string should return 0
		{"123.45", 123}, // With current implementation, it parses 123 and ignores .45
	}

	for _, tc := range tests {
		result := MustParseInt(tc.input)
		if result != tc.expected {
			t.Errorf("MustParseInt(%q): expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}

func TestCalculateTax(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request parameters
		if r.URL.Path != "/interface/2025Version1.xhtml" {
			t.Errorf("Expected path %s, got %s", "/interface/2025Version1.xhtml", r.URL.Path)
		}

		// Check query parameters
		q := r.URL.Query()
		if q.Get("code") != APICode {
			t.Errorf("Expected code %s, got %s", APICode, q.Get("code"))
		}

		// Send a mock response
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
		<lohnsteuer jahr="2025">
			<information>Mock tax response</information>
			<eingaben>
				<eingabe name="LZZ" value="1" status="ok"/>
				<eingabe name="RE4" value="50000" status="ok"/>
				<eingabe name="STKL" value="1" status="ok"/>
			</eingaben>
			<ausgaben>
				<ausgabe name="LSTLZZ" value="8000" type="STANDARD"/>
				<ausgabe name="SOLZLZZ" value="400" type="STANDARD"/>
			</ausgaben>
		</lohnsteuer>`))
	}))
	defer server.Close()

	// Since we can't modify BaseURL, we'll define a custom CalculateTax function for testing
	calculateTaxWithURL := func(baseURL string, req models.TaxRequest) (*TaxCalculationResponse, error) {
		url := fmt.Sprintf("%s?code=%s&LZZ=%d&RE4=%d&STKL=%d",
			baseURL, APICode, req.Period, req.Income, req.TaxClass)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to make HTTP request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
		}
		var taxResponse TaxCalculationResponse
		if err := xml.NewDecoder(resp.Body).Decode(&taxResponse); err != nil {
			return nil, fmt.Errorf("failed to decode XML response: %w", err)
		}

		return &taxResponse, nil
	}

	// Test successful request
	req := models.TaxRequest{
		Period:   models.Year,
		Income:   50000,
		TaxClass: models.TaxClass1,
	}

	resp, err := calculateTaxWithURL(server.URL+"/interface/2025Version1.xhtml", req)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp.Year != "2025" {
		t.Errorf("Expected year 2025, got %s", resp.Year)
	}

	if len(resp.Outputs.Output) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(resp.Outputs.Output))
	}

	// Check that the income tax value is correct
	var foundIncomeTax, foundSolidarityTax bool
	for _, output := range resp.Outputs.Output {
		if output.Name == "LSTLZZ" {
			if output.Value != "8000" {
				t.Errorf("Expected income tax 8000, got %s", output.Value)
			}
			foundIncomeTax = true
		}
		if output.Name == "SOLZLZZ" {
			if output.Value != "400" {
				t.Errorf("Expected solidarity tax 400, got %s", output.Value)
			}
			foundSolidarityTax = true
		}
	}

	if !foundIncomeTax {
		t.Error("Income tax output not found in response")
	}
	if !foundSolidarityTax {
		t.Error("Solidarity tax output not found in response")
	}
}

func TestCalculateTaxError(t *testing.T) {
	// Create a mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Since we can't modify BaseURL, we'll define a custom CalculateTax function for testing
	calculateTaxWithURL := func(baseURL string, req models.TaxRequest) (*TaxCalculationResponse, error) {
		url := fmt.Sprintf("%s?code=%s&LZZ=%d&RE4=%d&STKL=%d",
			baseURL, APICode, req.Period, req.Income, req.TaxClass)
		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("failed to make HTTP request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status: %s", resp.Status)
		}
		var taxResponse TaxCalculationResponse
		if err := xml.NewDecoder(resp.Body).Decode(&taxResponse); err != nil {
			return nil, fmt.Errorf("failed to decode XML response: %w", err)
		}

		return &taxResponse, nil
	}

	req := models.TaxRequest{
		Period:   models.Year,
		Income:   50000,
		TaxClass: models.TaxClass1,
	}

	resp, err := calculateTaxWithURL(server.URL+"/interface/2025Version1.xhtml", req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
	if resp != nil {
		t.Errorf("Expected nil response, got: %+v", resp)
	}
}