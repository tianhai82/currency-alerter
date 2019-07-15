package main

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

func makeGetRequest(url string, output interface{}) (err error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		err = errors.Wrap(err, "http get fails")
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		err = errors.Wrap(err, "json decoding fails")
		return
	}
	return
}

func makePostRequest(url string, input interface{}, output interface{}) (err error) {
	jsonValue, err := json.Marshal(input)
	if err != nil {
		err = errors.Wrap(err, "Unable to marshal input to json")
		return
	}
	resp, err := httpClient.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		err = errors.Wrap(err, "http post fails")
		return
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(output)
	if err != nil {
		err = errors.Wrap(err, "json decoding fails")
		return
	}
	return
}
