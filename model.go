package main

var AvailableCurrencies = []string{
	"CAD",
	"HKD",
	"ISK",
	"PHP",
	"DKK",
	"HUF",
	"CZK",
	"AUD",
	"RON",
	"SEK",
	"IDR",
	"INR",
	"BRL",
	"RUB",
	"HRK",
	"JPY",
	"THB",
	"CHF",
	"SGD",
	"PLN",
	"BGN",
	"TRY",
	"CNY",
	"NOK",
	"NZD",
	"ZAR",
	"USD",
	"MXN",
	"ILS",
	"GBP",
	"KRW",
	"MYR",
}

type ExchangeRatesResponse struct {
	Rates   map[string]map[string]float64 `json:"rates"`
	StartAt string                        `json:"start_at"`
	Base    string                        `json:"base"`
	EndAt   string                        `json:"end_at"`
}
type CurrencyHistory struct {
	TopCurrency  string
	BaseCurrency string
	Days         []string
	Rates        map[string]float64
}
type Subscription struct {
	UserID       int    `json:"userID"`
	TopCurrency  string `json:"topCurrency"`
	BaseCurrency string `json:"baseCurrency"`
}
type Msg struct {
	ChatID       int64  `json:"chat_id"`
	Text         string `json:"text"`
	ReplyToMsgID int    `json:"reply_to_message_id"`
}
type Alert struct {
}
