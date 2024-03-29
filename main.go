package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const ExchangeRateApi = "https://api.exchangeratesapi.io"

func main() {
	fmt.Println("Starting Currency Alerter")

	r := gin.Default()
	r.GET("/analyze", analyze)
	r.GET("/_ah/stop", shutdown)
	// r.GET("/setWebhook", setWebhook)
	r.POST("/telegramWebhook", webhookCallback)
	r.Run()
}
func shutdown(_ *gin.Context) {

}
func webhookCallback(c *gin.Context) {
	var telegramUpdate Update
	err := c.BindJSON(&telegramUpdate)
	if err != nil {
		fmt.Println(err)
		return
	}
	if telegramUpdate.Message.IsCommand() {
		err = handleUpdate(telegramUpdate)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func analyze(c *gin.Context) {
	fmt.Println("start analyze", time.Now())
	err := retrieveExchangeRates()
	if err != nil {
		fmt.Println(err)
		c.AbortWithError(http.StatusBadRequest, errors.New("fail to retrieve exchange rates"))
	}
}

// func setWebhook(_ *gin.Context) {
// 	setTelegramWebhook()
// }
