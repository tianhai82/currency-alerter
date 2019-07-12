package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/pkg/errors"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

const DateFormat = "2006-01-02"
const period = 30

func retrieveExchangeRates() (err error) {
	now := time.Now()
	endDate := now.Format(DateFormat)
	start := now.AddDate(0, 0, -period)
	startDate := start.Format(DateFormat)
	url := fmt.Sprintf("%s/history?start_at=%s&end_at=%s", ExchangeRateApi, startDate, endDate)
	fmt.Println(url)
	resp, err := httpClient.Get(url)
	if err != nil {
		err = errors.Wrap(err, "http get fails")
		return
	}
	defer resp.Body.Close()
	var target ExchangeRatesResponse
	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		err = errors.Wrap(err, "json decoding fails")
		return
	}
	histories := toCurrencyHistories(target)
	b, _ := json.Marshal(histories)
	fmt.Println(string(b))
	return
}
func toCurrencyHistories(rates ExchangeRatesResponse) (histories []CurrencyHistory) {
	base := rates.Base
	currencyHistories := make(map[string]CurrencyHistory)
	for day, currencies := range rates.Rates {
		for curr, rate := range currencies {
			dayRate := DayRate{
				Day:  day,
				Rate: rate,
			}
			currencyHistory, found := currencyHistories[curr]
			if found {
				rates := currencyHistory.Rates
				rates = append(rates, dayRate)
				currencyHistory.Rates = rates
			} else {
				currencyHistory.TopCurrency = curr
				currencyHistory.BaseCurrency = base
				rates := currencyHistory.Rates
				rates = append(rates, dayRate)
				currencyHistory.Rates = rates
			}
			currencyHistories[curr] = currencyHistory
		}
	}
	for _, history := range currencyHistories {
		sort.Slice(history.Rates, func(i, j int) bool {
			return history.Rates[i].Day > history.Rates[j].Day
		})
		histories = append(histories, history)
	}
	return
}
