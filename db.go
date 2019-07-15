package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
)

var client *firestore.Client

func init() {
	projectID := "currency-alerter"
	ctx := context.Background()
	var err error
	client, err = firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Println(err)
	}
}

func saveSubscription(userID int, topCurrency string, baseCurrency string) error {
	if client == nil {
		return errors.New("firestore client not initialised")
	}

	// TODO
	// convert currency to uppercase.
	// do not save duplicates
	//

	ctx := context.Background()
	_, _, err := client.Collection("subscription").Add(ctx, Subscription{
		UserID:       userID,
		TopCurrency:  topCurrency,
		BaseCurrency: baseCurrency,
	})
	return err
}
