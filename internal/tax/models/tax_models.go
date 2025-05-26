package models

type TaxClass int

const (
	TaxClass1 TaxClass = 1
	TaxClass2 TaxClass = 2
	TaxClass3 TaxClass = 3
	TaxClass4 TaxClass = 4
	TaxClass5 TaxClass = 5
	TaxClass6 TaxClass = 6
)

type PaymentPeriod int

const (
	Year   PaymentPeriod = 1
	Month  PaymentPeriod = 2
	Week   PaymentPeriod = 3
	Day    PaymentPeriod = 4
)

type TaxRequest struct {
	Period   PaymentPeriod
	Income   int
	TaxClass TaxClass
	
	AJAHR     int
	ALTER1    int
	KRV       int
	KVZ       float64
	PVS       int
	PVZ       int
	R         int
	ZKF       float64
	VBEZ      int
	VJAHR     int
	PKPV      int
	PKV       int
	PVA       int
}

type TaxResult struct {
	Income        float64
	IncomeTax     float64
	SolidarityTax float64
	TotalTax      float64
	NetIncome     float64
	TaxRate       float64
	Error         error
}