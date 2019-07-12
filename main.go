package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const ExchangeRateApi = "https://api.exchangeratesapi.io"

func main() {
	fmt.Println("Starting Currency Alerter")

	r := gin.Default()
	r.GET("/analyze", analyze)

	r.Run()
}

func analyze(_ *gin.Context) {
	retrieveExchangeRates()
}
