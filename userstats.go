package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Stats struct {
	OneNameErrorResponse
	Registrations int `json:"registrations"`
}

func (s Stats) String() string {
	if s.Error != nil {
		return s.Error.String()
	}
	return "registrations:" + fmt.Sprintln(s.Registrations)
}

func (c *client) GetUserStats() (Stats, error) {
	url := c.BaseUrl + "/stats/users"
	var jsonObj Stats
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
