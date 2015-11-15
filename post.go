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
type OneNamePost struct {
	Passname         string `json:"passname,omitempty"`
	RecipientAddress string `json:"recipient_address,omitempty"`
	SignedHex        string `json:"signed_hex,omitempty"`
}

func (c *client) PostRequest(url string, payload OneNamePost) (OneNameErrorResponse, error) {
	responseObject := OneNameErrorResponse{}

	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return responseObject, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return responseObject, err
	}

	req.SetBasicAuth(c.ApiID, c.ApiSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return responseObject, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return responseObject, err
	}

	x := &responseObject
	err = json.Unmarshal(body, x)
	if err != nil {
		return responseObject, err
	}

	if x.Error != nil {
		err = errors.New("Error: " + x.Error.Type + " - " + x.Error.Message)
	}

	return responseObject, err
}

func (c *client) RegisterUser(username, address string) (OneNameErrorResponse, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return OneNameErrorResponse{}, err
	}
	url := c.BaseUrl + "/users"
	payload := OneNamePost{username, address, ""}
	return c.PostRequest(url, payload)
}

func (c *client) BroadcastTransactions(signedHex string) (OneNameErrorResponse, error) {
	url := c.BaseUrl + "/transactions"
	payload := OneNamePost{"", "", signedHex}
	return c.PostRequest(url, payload)
}
