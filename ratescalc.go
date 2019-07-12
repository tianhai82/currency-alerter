package main

func monitorRates(histories []CurrencyHistory) {
	// loop thru each currency pair
	// calculate std dev, max dev and (max + std) / 2
	// if current deviation with mean is more than (max+std)/2, send alert to subscribed users
}
