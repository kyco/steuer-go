package bmf

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"tax-calculator/internal/tax/models"
	"testing"
)

func TestMustParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"100", 100},
		{"12345", 12345},
		{"-999", -999},
		{"", 0},
		{"abc", 0},
	}

	for _, tc := range tests {
		result := MustParseInt(tc.input)
		if result != tc.expected {
			t.Errorf("MustParseInt(%q): expected %d, got %d", tc.input, tc.expected, result)
		}
	}
}

// Custom HTTP client for testing
type mockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, fmt.Errorf("no mock implementation")
}

// Helper function to test CalculateTax with a custom client
func calculateTaxWithClient(client *http.Client, baseURL string, req models.TaxRequest) (*TaxCalculationResponse, error) {
	url := fmt.Sprintf("%s?code=%s&LZZ=%d&RE4=%d&STKL=%d",
		baseURL, APICode, req.Period, req.Income, req.TaxClass)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := client.Do(httpReq)
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

func TestCalculateTax(t *testing.T) {
	// Create a test server that returns a mock XML response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request parameters
		query := r.URL.Query()
		if code := query.Get("code"); code != APICode {
			t.Errorf("Expected code=%s, got %s", APICode, code)
		}

		lzz := query.Get("LZZ")
		re4 := query.Get("RE4")
		stkl := query.Get("STKL")

		if lzz == "" || re4 == "" || stkl == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Return a mock XML response
		xmlResponse := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<lohnsteuer jahr="2025">
			<information>Mock tax response</information>
			<eingaben>
				<eingabe name="LZZ" value="%s" status="ok"/>
				<eingabe name="RE4" value="%s" status="ok"/>
				<eingabe name="STKL" value="%s" status="ok"/>
			</eingaben>
			<ausgaben>
				<ausgabe name="LSTLZZ" value="800000" type="STANDARD"/>
				<ausgabe name="SOLZLZZ" value="40000" type="STANDARD"/>
			</ausgaben>
		</lohnsteuer>`, lzz, re4, stkl)

		w.Header().Set("Content-Type", "application/xml")
		fmt.Fprint(w, xmlResponse)
	}))
	defer server.Close()

	// Create a custom client for testing
	client := &http.Client{}

	// Test the CalculateTax function with our custom URL
	req := models.TaxRequest{
		Period:   models.Year,
		Income:   5000000, // 50,000.00 EUR in cents
		TaxClass: models.TaxClass1,
	}

	// Call the function with our test server URL
	result, err := calculateTaxWithClient(client, server.URL, req)
	if err != nil {
		t.Fatalf("CalculateTax failed: %v", err)
	}

	// Check the result
	if result.Year != "2025" {
		t.Errorf("Expected year 2025, got %s", result.Year)
	}

	if result.Information != "Mock tax response" {
		t.Errorf("Expected 'Mock tax response', got %s", result.Information)
	}

	// Check inputs
	if len(result.Inputs.Input) != 3 {
		t.Errorf("Expected 3 inputs, got %d", len(result.Inputs.Input))
	}

	// Check outputs
	if len(result.Outputs.Output) != 2 {
		t.Errorf("Expected 2 outputs, got %d", len(result.Outputs.Output))
	}

	var incomeTax, solidarityTax string
	for _, output := range result.Outputs.Output {
		if output.Name == "LSTLZZ" {
			incomeTax = output.Value
		} else if output.Name == "SOLZLZZ" {
			solidarityTax = output.Value
		}
	}

	if incomeTax != "800000" {
		t.Errorf("Expected income tax 800000, got %s", incomeTax)
	}

	if solidarityTax != "40000" {
		t.Errorf("Expected solidarity tax 40000, got %s", solidarityTax)
	}
}

func TestCalculateTaxServerError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	// Create a custom client for testing
	client := &http.Client{}

	// Test the CalculateTax function
	req := models.TaxRequest{
		Period:   models.Year,
		Income:   5000000,
		TaxClass: models.TaxClass1,
	}

	// Call the function
	_, err := calculateTaxWithClient(client, server.URL, req)

	// We should get an error
	if err == nil {
		t.Error("Expected error on server error, got nil")
	}

	if !strings.Contains(err.Error(), "API request failed with status") {
		t.Errorf("Expected error about API status, got: %v", err)
	}
}
