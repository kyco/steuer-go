package models

// TaxClass (Steuerklasse) - German tax classes
type TaxClass int

const (
	TaxClass1 TaxClass = 1
	TaxClass2 TaxClass = 2
	TaxClass3 TaxClass = 3
	TaxClass4 TaxClass = 4
	TaxClass5 TaxClass = 5
	TaxClass6 TaxClass = 6
)

// PaymentPeriod (Lohnzahlungszeitraum - LZZ)
type PaymentPeriod int

const (
	Year   PaymentPeriod = 1
	Month  PaymentPeriod = 2
	Week   PaymentPeriod = 3
	Day    PaymentPeriod = 4
)

// TaxRequest represents a request to calculate tax
type TaxRequest struct {
	Period   PaymentPeriod // LZZ
	Income   int           // RE4 (in cents)
	TaxClass TaxClass      // STKL
}

// TaxResult contains summarized tax calculation results
type TaxResult struct {
	Income        float64
	IncomeTax     float64
	SolidarityTax float64
	TotalTax      float64
	NetIncome     float64
	TaxRate       float64
	Error         error
}