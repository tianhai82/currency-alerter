package main

import (
	"fmt"
	"math"

	"github.com/pkg/errors"
)

const tianhaiID = 21450012

var ratesCache = make(map[string]CurrencyHistory)

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
		alert := getAlert(curHist, sub.ChatID)
		if alert.TopCurrency != "" {
			sendAlert(alert)
			if alert.ChatID == tianhaiID {
				errSaveOrder := saveOrder(alert)
				if errSaveOrder != nil {
					fmt.Println(errSaveOrder)
				}
			}
		}
	}
	return nil
}

func getAlert(curHist CurrencyHistory, userID int64) Alert {
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
	currentRate := curHist.Rates[curHist.Days[0]]
	currentDev := currentRate - mean
	if math.Abs(currentDev) >= limit {
		fmt.Printf("sending alert for %s/%s\n", curHist.TopCurrency, curHist.BaseCurrency)
		return Alert{
			Subscription: Subscription{
				ChatID:       userID,
				TopCurrency:  curHist.TopCurrency,
				BaseCurrency: curHist.BaseCurrency,
			},
			CurrentDevPercen: currentDev / currentRate * 100,
			StdDevPercen:     stdDev / currentRate * 100,
			MaxDevPercen:     maxDev / currentRate * 100,
			CurrentRate:      currentRate,
			MeanRate:         mean,
			LimitDevPercen:   limit / currentRate * 100,
		}
	}
	fmt.Printf("no alert for %s/%s.\n", curHist.TopCurrency, curHist.BaseCurrency)
	return Alert{}
}
func sendAlert(alert Alert) {
	rec := ""
	if alert.CurrentDevPercen < 0 {
		rec = "SELL"
	} else {
		rec = "BUY"
	}
	text := fmt.Sprintf(`__Alert for *%s/%s*__
Current Rate: *%.5f*
Moving Avg Rate: *%.5f*
Standard Dev: *%.6f%%*
Max Dev: *%.6f%%*
Limit Dev: *%.6f%%*
Current Dev: *%.6f%%*

Recommendation: *%s*`,
		alert.TopCurrency, alert.BaseCurrency,
		alert.CurrentRate, alert.MeanRate,
		alert.StdDevPercen, alert.MaxDevPercen,
		alert.LimitDevPercen, alert.CurrentDevPercen,
		rec)
	sendMessage(Msg{
		ChatID:    alert.ChatID,
		Text:      text,
		ParseMode: "Markdown",
	})
}

func getPairRates(sub Subscription, histories map[string]CurrencyHistory) (CurrencyHistory, error) {
	// if found in cache, return
	key := fmt.Sprintf("%s/%s", sub.TopCurrency, sub.BaseCurrency)
	if rate, found := ratesCache[key]; found {
		return rate, nil
	}

	// if top currency is EUR, invert it's base
	if sub.TopCurrency == "EUR" {
		cur, found := histories[sub.BaseCurrency]
		if !found {
			return CurrencyHistory{}, errors.New(sub.TopCurrency + " not found")
		}
		history := CurrencyHistory{
			TopCurrency:  sub.TopCurrency,
			BaseCurrency: sub.BaseCurrency,
			Rates:        make(map[string]float64),
		}
		for _, day := range cur.Days {
			rate, foundRate := cur.Rates[day]
			if !foundRate {
				return CurrencyHistory{}, errors.Errorf("%s not found in %s history", day, sub.BaseCurrency)
			}
			history.Rates[day] = 1.0 / rate
			history.Days = append(history.Days, day)
		}
		ratesCache[key] = history // save to cache before returning
		return history, nil
	}

	// if top currency not found, return error
	topCur, found := histories[sub.TopCurrency]
	if !found {
		return CurrencyHistory{}, errors.New(sub.TopCurrency + " not found")
	}

	// if base is EUR, return top currency
	if sub.BaseCurrency == "EUR" {
		return topCur, nil
	}

	// if base currency not found, return error
	baseCur, found := histories[sub.BaseCurrency]
	if !found {
		return CurrencyHistory{}, errors.New(sub.BaseCurrency + " not found")
	}

	// if both found, divide the rates for each day
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
	ratesCache[key] = history // save to cache before returning
	return history, nil
}
