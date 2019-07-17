package main

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

func monitorRates(histories map[string]CurrencyHistory) error {
	subscriptions, err := retrieveAllSubscriptions()
	if err != nil {
		return err
	}
	for _, sub := range subscriptions {
		curHist, err := getPairRates(sub, histories)
		if err != nil {
			fmt.Println(err)
			continue
		}
		alert, err := getAlert(curHist, sub.ChatID)
		if err == nil {
			sendAlert(alert)
		} else {
			fmt.Println(err)
		}
	}
	return nil

}

func getAlert(curHist CurrencyHistory, userID int64) (Alert, error) {
	// calculate std dev, max dev and (max + std) / 2
	// if current deviation with mean is more than (max+std)/2, send alert to subscribed users
	maxDev := 0.0
	total := 0.0
	for _, day := range curHist.Days {
		total += curHist.Rates[day]
	}
	mean := total / float64(len(curHist.Days))
	varianceSum := 0.0
	for _, day := range curHist.Days {
		rate := curHist.Rates[day]
		diff := math.Abs(rate - mean)
		if diff > maxDev {
			maxDev = diff
		}
		varianceSum += math.Pow(diff, 2.0)
	}
	variance := varianceSum / float64(len(curHist.Days))
	stdDev := math.Sqrt(variance)
	limit := (maxDev + stdDev) / 2
	currentDev := math.Abs(curHist.Rates[curHist.Days[0]] - mean)
	if currentDev >= limit {
		return Alert{}, nil
	}
	currentRate := curHist.Rates[curHist.Days[0]]
	text := fmt.Sprintf("Std Dev: %.6f%%. Max Dev: %.6f%%. Current Dev: %.6f%%. Current rate %s/%s: %.5f",
		(stdDev/currentRate)*100, (maxDev/currentRate)*100, (currentDev/currentRate)*100,
		curHist.TopCurrency, curHist.BaseCurrency, currentRate)
	sendMessage(Msg{
		ChatID: int64(userID),
		Text:   text,
	})
	return Alert{}, errors.New("current deviation within limit")
}
func sendAlert(alert Alert) {

}

func getPairRates(sub Subscription, histories map[string]CurrencyHistory) (CurrencyHistory, error) {
	topCur, found := histories[sub.TopCurrency]
	if !found {
		return CurrencyHistory{}, errors.New(sub.TopCurrency + " not found")
	}
	if sub.BaseCurrency == "EUR" {
		return topCur, nil
	}
	baseCur, found := histories[sub.BaseCurrency]
	if !found {
		return CurrencyHistory{}, errors.New(sub.BaseCurrency + " not found")
	}
	history := CurrencyHistory{
		TopCurrency:  sub.TopCurrency,
		BaseCurrency: sub.BaseCurrency,
		Rates:        make(map[string]float64),
	}
	for _, day := range topCur.Days {
		topRate, found := topCur.Rates[day]
		if !found {
			return CurrencyHistory{}, errors.Errorf("%s not found in %s history", day, sub.TopCurrency)
		}
		baseRate, found := baseCur.Rates[day]
		if !found {
			return CurrencyHistory{}, errors.Errorf("%s not found in %s history", day, sub.BaseCurrency)
		}
		rate := topRate / baseRate
		history.Rates[day] = rate
		history.Days = append(history.Days, day)
	}
	return history, nil
}
