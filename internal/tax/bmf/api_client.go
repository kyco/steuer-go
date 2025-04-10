package bmf

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"tax-calculator/internal/tax/models"
)

const (
	BaseURL = "http://www.bmf-steuerrechner.de/interface/2025Version1.xhtml"
	APICode = "extS2025"
)

type TaxCalculationResponse struct {
	XMLName     xml.Name `xml:"lohnsteuer"`
	Year        string   `xml:"jahr,attr"`
	Information string   `xml:"information"`
	Inputs      Inputs   `xml:"eingaben"`
	Outputs     Outputs  `xml:"ausgaben"`
}

type Inputs struct {
	Input []Input `xml:"eingabe"`
}

type Input struct {
	Name   string `xml:"name,attr"`
	Value  string `xml:"value,attr"`
	Status string `xml:"status,attr"`
}

type Outputs struct {
	Output []Output `xml:"ausgabe"`
}

type Output struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Type  string `xml:"type,attr"`
}

func CalculateTax(req models.TaxRequest) (*TaxCalculationResponse, error) {
	url := fmt.Sprintf("%s?code=%s&LZZ=%d&RE4=%d&STKL=%d",
		BaseURL, APICode, req.Period, req.Income, req.TaxClass)
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

func MustParseInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0
	}
	return result
}