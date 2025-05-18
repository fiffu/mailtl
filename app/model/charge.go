package model

type Charge struct {
	LocalCurrency  string
	LocalAmount    float64
	CardNumber     string
	Timestamp      string
	ChargeCurrency string
	ChargeAmount   float64
	Purpose        string
}
