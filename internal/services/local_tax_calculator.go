package services

import (
	"fmt"
	"math/big"
	"sync"
	"tax-calculator/internal/adapters/api"
	"tax-calculator/internal/domain/models"
)

type LocalTaxCalculator struct {
	xmlData     *api.PAPData
	calculator  *api.TaxCalculator
	initialized bool
	mu          sync.RWMutex
}

var (
	localCalculator *LocalTaxCalculator
	once            sync.Once
)

func GetLocalTaxCalculator() *LocalTaxCalculator {
	once.Do(func() {
		localCalculator = &LocalTaxCalculator{
			initialized: false,
		}
	})
	return localCalculator
}

func (l *LocalTaxCalculator) Initialize() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.initialized {
		return nil
	}

	xmlData, err := api.FetchTaxCalculationXML()
	if err != nil {
		return fmt.Errorf("failed to initialize local tax calculator: %w", err)
	}

	l.xmlData = xmlData
	l.calculator = api.NewTaxCalculator(xmlData)
	l.initialized = true

	return nil
}

func (l *LocalTaxCalculator) IsInitialized() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.initialized
}

func (l *LocalTaxCalculator) CalculateTax(req models.TaxRequest) (*api.TaxCalculationResponse, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.initialized {
		return nil, fmt.Errorf("local tax calculator not initialized")
	}

	l.calculator.SetInputValue("LZZ", int(req.Period))
	l.calculator.SetInputValue("RE4", req.Income)
	l.calculator.SetInputValue("STKL", int(req.TaxClass))

	l.calculator.SetInputValue("R", req.R)
	l.calculator.SetInputValue("AJAHR", req.AJAHR)
	l.calculator.SetInputValue("ALTER1", req.ALTER1)
	l.calculator.SetInputValue("KRV", req.KRV)
	l.calculator.SetInputValue("KVZ", req.KVZ)
	l.calculator.SetInputValue("PVS", req.PVS)
	l.calculator.SetInputValue("PVZ", req.PVZ)
	l.calculator.SetInputValue("PKV", req.PKV)
	l.calculator.SetInputValue("PVA", req.PVA)
	l.calculator.SetInputValue("ZKF", req.ZKF)
	l.calculator.SetInputValue("VBEZ", int(req.VBEZ * 100))
	l.calculator.SetInputValue("VJAHR", req.VJAHR)
	l.calculator.SetInputValue("PKPV", int(req.PKPV * 100))
	
	l.calculator.SetInputValue("JFREIB", 0)
	l.calculator.SetInputValue("JHINZU", 0)
	l.calculator.SetInputValue("JRE4", 0)
	l.calculator.SetInputValue("JRE4ENT", 0)
	l.calculator.SetInputValue("JVBEZ", 0)
	l.calculator.SetInputValue("VBEZM", 0)
	l.calculator.SetInputValue("VBEZS", 0)
	l.calculator.SetInputValue("VBS", 0)
	l.calculator.SetInputValue("VKAPA", 0)
	l.calculator.SetInputValue("VMT", 0)
	l.calculator.SetInputValue("ZMVB", 12)
	l.calculator.SetInputValue("SONSTB", 0)
	l.calculator.SetInputValue("SONSTENT", 0)
	l.calculator.SetInputValue("STERBE", 0)
	l.calculator.SetInputValue("af", 0)
	l.calculator.SetInputValue("f", 1.0)

	if err := l.calculator.Calculate(); err != nil {
		return nil, fmt.Errorf("tax calculation failed: %w", err)
	}

	response := &api.TaxCalculationResponse{
		Year:        "2025",
		Information: "Local calculation based on BMF XML",
		Outputs: api.Outputs{
			Output: make([]api.Output, 0),
		},
	}

	for name, value := range l.calculator.OutputValues {
		var strValue string
		switch v := value.(type) {
		case int:
			strValue = fmt.Sprintf("%d", v)
		case *big.Rat:
			// Convert to cents integer
			cents := new(big.Rat).Mul(v, big.NewRat(100, 1))
			// Convert to integer string
			f, _ := cents.Float64()
			strValue = fmt.Sprintf("%d", int(f))
		default:
			strValue = fmt.Sprintf("%v", v)
		}

		response.Outputs.Output = append(response.Outputs.Output, api.Output{
			Name:  name,
			Value: strValue,
			Type:  "BigDecimal",
		})
	}

	return response, nil
}