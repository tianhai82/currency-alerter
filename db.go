package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"github.com/pkg/errors"
	"google.golang.org/api/iterator"
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

func retrieveAllSubscriptions() ([]Subscription, error) {
	if client == nil {
		return nil, errors.New("firestore client not initialised")
	}

	ctx := context.Background()
	iter := client.Collection("subscription").Documents(ctx)
	defer iter.Stop()
	var subs []Subscription
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "fail to retrieve all subscriptions")
		}
		var sub Subscription
		errData := doc.DataTo(&sub)
		if errData != nil {
			fmt.Println(errData)
			continue
		} else {
			subs = append(subs, sub)
		}
	}
	return subs, nil
}

func saveSubscription(userID int, topCurrency string, baseCurrency string) error {
	if client == nil {
		return errors.New("firestore client not initialised")
	}

	ctx := context.Background()
	key := fmt.Sprintf("%d-%s-%s", userID, topCurrency, baseCurrency)
	_, err := client.Collection("subscription").Doc(key).Set(ctx, Subscription{
		UserID:       userID,
		TopCurrency:  topCurrency,
		BaseCurrency: baseCurrency,
	})
	return err
}
