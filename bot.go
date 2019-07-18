package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const botToken = "801803446:AAFgJE90u7rYVe4tHJ-fEM_tNau3pWXUFBg"

const botURL = "https://api.telegram.org/bot"

const webhookURL = "https://currency-alerter.appspot.com/telegramWebhook"

type webhook struct {
	URL            string   `json:"url"`
	MaxConnections uint     `json:"max_connections"`
	AllowedUpdates []string `json:"allowed_updates"`
}

func handleUpdate(update Update) error {
	if !update.Message.IsCommand() {
		return errors.New("Bot only handles command")
	}
	cmd := update.Message.Command()
	switch cmd {
	case "sub":
		return subscribe(update.Message)
	case "start":
		return sendMessage(Msg{
			ChatID: update.Message.Chat.ID,
			Text:   "Subscribe to currency pair by entering the '/sub' command",
		})
	default:
		return errors.New("Invalid command")
	}
}

func sendMessage(msg Msg) error {
	url := fmt.Sprintf("%s%s/sendMessage", botURL, botToken)
	var message Message
	return makePostRequest(url, msg, &message)
}

func subscribe(msg *Message) error {
	args := msg.CommandArguments()
	currencies := strings.Split(args, "/")
	for i, s := range currencies {
		currencies[i] = strings.ToUpper(s)
	}
	if err := checkCurrenciesArgs(currencies); err != nil {
		sendMessage(Msg{
			ChatID:       msg.Chat.ID,
			ReplyToMsgID: msg.MessageID,
			Text:         err.Error(),
		})
		return err
	}
	err := saveSubscription(msg.Chat.ID, currencies[0], currencies[1])
	if err != nil {
		sendMessage(Msg{
			ChatID:       msg.Chat.ID,
			ReplyToMsgID: msg.MessageID,
			Text:         "Unsuccessul. Error saving subscription. Please try again later",
		})
		return err
	}
	sendMessage(Msg{
		ChatID:       msg.Chat.ID,
		ReplyToMsgID: msg.MessageID,
		Text:         fmt.Sprintf("Subscribed to %s/%s successfully!", strings.ToUpper(currencies[0]), strings.ToUpper(currencies[1])),
	})
	return nil
}

func checkCurrenciesArgs(currencies []string) error {
	if len(currencies) != 2 {
		return errors.New("Unsuccessful. Please provide Currency_Code_1/Currency_Code_2")
	}

	if !contains(AvailableCurrencies, currencies[0]) || !contains(AvailableCurrencies, currencies[1]) {
		return errors.Errorf("Available currencies: %s", strings.Join(AvailableCurrencies, ", "))
	}

	if currencies[0] == currencies[1] {
		return errors.New("Currency 1 cannot be the same as currency 2")
	}
	return nil
}

// func setTelegramWebhook() error {
// 	url := fmt.Sprintf("%s%s/setWebhook", botUrl, botToken)
// 	fmt.Println(url)
// 	var out interface{}
// 	err := makePostRequest(url, webhook{
// 		Url: webhookURL,
// 	}, &out)
// 	if err != nil {
// 		fmt.Println(err)
// 		return errors.Wrap(err, "fail to create webhook")
// 	}
// 	b, _ := json.Marshal(out)
// 	fmt.Println(string(b))
// 	return nil
// }
