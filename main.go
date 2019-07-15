package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

const ExchangeRateApi = "https://api.exchangeratesapi.io"

func main() {
	fmt.Println("Starting Currency Alerter")

	r := gin.Default()
	r.GET("/analyze", analyze)
	r.GET("/setWebhook", setWebhook)
	r.POST("/telegramWebhook", webhookCallback)
	r.Run()
}

func webhookCallback(c *gin.Context) {
	var telegramMsg map[string]interface{}
	err := c.BindJSON(&telegramMsg)
	if err != nil {
		fmt.Println(err)
	}
	b, _ := json.Marshal(telegramMsg)
	fmt.Println(string(b))
}

func analyze(_ *gin.Context) {
	retrieveExchangeRates()
}
func setWebhook(_ *gin.Context) {
	setTelegramWebhook()
}
