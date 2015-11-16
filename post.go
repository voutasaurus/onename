package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/btcsuite/btcutil/base58"
)

// POST

type RegisterUserPost struct {
	Passname         string `json:"passname,omitempty"`
	RecipientAddress string `json:"recipient_address,omitempty"`
}

type BroadcastTransactionPost struct {
	SignedHex string `json:"signed_hex,omitempty"`
}

// Best approximation of an OR type for RegisterUserPost and BroadcastTransactionPost
type OneNamePost struct {
	RegisterUserPost
	BroadcastTransactionPost
}

func (c *client) PostRequest(url string, payload OneNamePost) (OneNameErrorResponse, error) {
	responseObject := OneNameErrorResponse{}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return responseObject, errors.New("Request not sent. Could not JSONify payload: " + err.Error())
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return responseObject, errors.New("Request not sent. Could not construct request.\n\t Check the URL \"" + url + "\" - got error: " + err.Error())
	}

	req.SetBasicAuth(c.ApiID, c.ApiSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseObject, errors.New("Request failed.\n\t Got error: " + err.Error())
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseObject, errors.New("Request completed but read of response body failed.\n\t Got error: " + err.Error())
	}

	jsonObj := &responseObject
	err = json.Unmarshal(body, jsonObj)
	if err != nil {
		return *jsonObj, JSONRead(err, body)
	}

	if jsonObj.Error != nil {
		err = errors.New("Error: " + jsonObj.Error.Type + " - " + jsonObj.Error.Message)
	}

	return responseObject, err
}

func (c *client) RegisterUser(username, address string) (OneNameErrorResponse, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return OneNameErrorResponse{}, errors.New("Request not sent. The address provided is not a valid bitcoin address: " + err.Error())
	}
	url := c.BaseUrl + "/users"
	payload := OneNamePost{RegisterUserPost{username, address}, BroadcastTransactionPost{}}
	return c.PostRequest(url, payload)
}

func (c *client) BroadcastTransactions(signedHex string) (OneNameErrorResponse, error) {
	url := c.BaseUrl + "/transactions"
	payload := OneNamePost{RegisterUserPost{}, BroadcastTransactionPost{signedHex}}
	return c.PostRequest(url, payload)
}
