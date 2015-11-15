package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/base58"
)

type UnspentsResponse struct {
	OneNameErrorResponse
	Unspents []Unspents `json:"unspents"`
}

func (s UnspentsResponse) String() string {
	if s.Error != nil {
		return s.Error.String()
	}
	return "Unspent: " + fmt.Sprintln(s.Unspents)
}

type Unspents struct {
	Confirmations   int    `json:"confirmations"`
	OutputIndex     int    `json:"output_index"`
	ScriptHex       string `json:"script_hex"`
	ScriptOpCodes   string `json:"script_opcodes"`
	ScriptType      string `json:"pubkeyhash"`
	TransactionHash string `json:"transaction_hash"`
	Value           int    `json:"value"`
}

func (c *client) GetUnspents(address string) (UnspentsResponse, error) {
	var jsonObj UnspentsResponse
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return jsonObj, errors.New("Address must be a valid cryptocurrency address. Instead it is " + address + err.Error())
	}
	url := c.BaseUrl + "/addresses/" + address + "/unspents"
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
