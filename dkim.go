package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type DKIMResponse struct {
	OneNameErrorResponse
	KeyCurve  string `json:"key_curve"`
	KeyType   string `json:"key_type"`
	PublicKey string `json:"public_key"`
}

func (s DKIMResponse) String() string {
	if s.Error != nil {
		return s.Error.String()
	}
	ret := "DKIM:\n"
	ret += fmt.Sprintln("\tcurve:", s.KeyCurve)
	ret += fmt.Sprintln("\tkey type:", s.KeyType)
	ret += fmt.Sprintln("\tkey type:", s.PublicKey)
	return ret
}

func (c *client) GetDKIMInfo(domain string) (DKIMResponse, error) {
	var jsonObj DKIMResponse
	url := c.BaseUrl + "/domains/" + domain + "/dkim"
	jsonBytes, err := c.GetRequest(url)
	if err != nil {
		return jsonObj, err
	}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return jsonObj, JSONRead(err, jsonBytes)
	}
	if jsonObj.Error != nil {
		err = errors.New("Error: " + jsonObj.Error.Type + " - " + jsonObj.Error.Message)
	}
	return jsonObj, err
}
