package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

type NamesResponse struct {
	OneNameErrorResponse
	Names []string `json:"names"`
}

func (s NamesResponse) String() string {
	if s.Error != nil {
		return s.Error.String()
	}
	return "Names: " + fmt.Sprintln(s.Names)
}

func (c *client) GetNames(address string) (NamesResponse, error) {
	var jsonObj NamesResponse
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return jsonObj, errors.New("Address must be a valid cryptocurrency address. Instead it is " + address + err.Error())
	}
	url := c.BaseUrl + "/addresses/" + address + "/names"
	jsonBytes, err := c.GetRequest(url)
	if err != nil {
		return jsonObj, err
	}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return jsonObj, err
	}
	if jsonObj.Error != nil {
		err = errors.New("Error: " + jsonObj.Error.Type + " - " + jsonObj.Error.Message)
	}
	return jsonObj, err
}
