package main

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

const botToken = "801803446:AAFgJE90u7rYVe4tHJ-fEM_tNau3pWXUFBg"

const botUrl = "https://api.telegram.org/bot"

const webhookURL = "https://currency-alerter.appspot.com/telegramWebhook"

type webhook struct {
	Url            string   `json:"url"`
	MaxConnections uint     `json:"max_connections"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func setTelegramWebhook() error {
	url := fmt.Sprintf("%s%s/setWebhook", botUrl, botToken)
	fmt.Println(url)
	var out interface{}
	err := makePostRequest(url, webhook{
		Url: webhookURL,
	}, &out)
	if err != nil {
		fmt.Println(err)
		return errors.Wrap(err, "fail to create webhook")
	}
	b, _ := json.Marshal(out)
	fmt.Println(string(b))
	return nil
}
