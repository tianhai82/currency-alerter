package main

import (
	"time"
)

var AvailableCurrencies = []string{
	"EUR",
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
	ChatID       int64  `json:"chatID"`
	TopCurrency  string `json:"topCurrency"`
	BaseCurrency string `json:"baseCurrency"`
}
type Msg struct {
	ChatID       int64  `json:"chat_id"`
	Text         string `json:"text"`
	ReplyToMsgID int    `json:"reply_to_message_id"`
	ParseMode    string `json:"parse_mode"`
}
type Alert struct {
	Subscription     `json:",inline"`
	CurrentDevPercen float64 `json:"currentDevPercen"`
	StdDevPercen     float64 `json:"stdDevPercen"`
	MaxDevPercen     float64 `json:"maxDevPercen"`
	LimitDevPercen   float64 `json:"limitDevPercen"`
	CurrentRate      float64 `json:"currentRate"`
	MeanRate         float64 `json:"meanRate"`
}
type Order struct {
	BuyOrSell    string    `json:"buyOrSell"`
	Price        float64   `json:"price"`
	TopCurrency  string    `json:"topCurrency"`
	BaseCurrency string    `json:"baseCurrency"`
	DateTime     time.Time `json:"dateTime"`
}
