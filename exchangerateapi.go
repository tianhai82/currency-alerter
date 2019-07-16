package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"
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
	var target ExchangeRatesResponse
	err = makeGetRequest(url, &target)
	if err != nil {
		return
	}
	histories := toCurrencyHistories(target)
	monitorRates(histories)
	return
}
func toCurrencyHistories(rates ExchangeRatesResponse) (currencyHistories map[string]CurrencyHistory) {
	base := rates.Base
	currencyHistories = make(map[string]CurrencyHistory)
	for day, currencies := range rates.Rates {
		for curr, rate := range currencies {
			currencyHistory, found := currencyHistories[curr]
			if found {
				currencyHistory.Rates[day] = rate
				days := currencyHistory.Days
				days = append(days, day)
				currencyHistory.Days = days
			} else {
				currencyHistory.TopCurrency = curr
				currencyHistory.BaseCurrency = base
				currencyHistory.Rates = make(map[string]float64)
				currencyHistory.Rates[day] = rate
				days := currencyHistory.Days
				days = append(days, day)
				currencyHistory.Days = days
			}
			currencyHistories[curr] = currencyHistory
		}
	}
	for _, history := range currencyHistories {
		sort.Slice(history.Days, func(i, j int) bool {
			return history.Days[i] > history.Days[j]
		})
	}
	return
}
