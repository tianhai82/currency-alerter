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
	Days         []string
	Rates        map[string]float64
}
type Subscription struct {
	UserID       int
	TopCurrency  string
	BaseCurrency string
}
type Msg struct {
	ChatID       int64  `json:"chat_id"`
	Text         string `json:"text"`
	ReplyToMsgID int    `json:"reply_to_message_id"`
}
