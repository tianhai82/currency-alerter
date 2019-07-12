package main

type ExchangeRatesResponse struct {
	Rates   map[string]map[string]float64 `json:"rates"`
	StartAt string                        `json:"start_at"`
	Base    string                        `json:"base"`
	EndAt   string                        `json:"end_at"`
}
type CurrencyHistory struct {
	TopCurrency  string
	BaseCurrency string
	Rates        []DayRate
}
type DayRate struct {
	Day  string
	Rate float64
}
