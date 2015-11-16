package main

import (
	"encoding/json"
	"errors"
)

type AllUsersResponse struct {
	OneNameErrorResponse
	Stats     Stats              `json:"stats"`
	Usernames []string           `json:"usernames"`
	Profiles  map[string]Profile `json:"profiles"`
}

// GetAllUsers
func (c *client) GetAllUsers() (AllUsersResponse, error) {
	var jsonObj AllUsersResponse
	url := c.BaseUrl + "/users"
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
