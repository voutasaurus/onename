package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type LookupResponse map[string]User

type User struct {
	Profile       Profile
	Verifications []Verification
}

func (u User) String() string {
	ret := fmt.Sprint(u.Profile.String())
	ret += "\n\tVerifications: "
	ret += fmt.Sprint(u.Verifications)
	return ret
}

type Verification struct {
	Identifier string `json:"identifier"`
	ProofURL   string `json:"proof_url"`
	Service    string `json:"service"`
	Valid      bool   `json:"valid"`
}

func (v Verification) String() string {
	ret := "\n"
	ret += "\t\tidentifier: " + v.Identifier + "\n"
	ret += "\t\tproof_url: " + v.ProofURL + "\n"
	ret += "\t\tservice: " + v.Service + "\n"
	ret += "\t\tvalid: " + fmt.Sprint(v.Valid) + "\n"
	return ret
}

// GetUsers Looks up users in the input slice
func (c *client) GetUsers(usernames []string) (LookupResponse, error) {
	url := c.BaseUrl + "/users/" + strings.Join(usernames, ",")
	return c.GetUserObjects(url)
}

// GetUserObject parses the JSON responses matching the LookupResponse type
func (c *client) GetUserObjects(url string) (LookupResponse, error) {
	var jsonObj LookupResponse
	jsonBytes, err := c.GetRequest(url)
	if err != nil {
		return jsonObj, err
	}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return jsonObj, JSONRead(err, jsonBytes)
	}
	if _, errorPresent := jsonObj["error"]; errorPresent {
		var errResp OneNameErrorResponse
		err = json.Unmarshal(jsonBytes, &errResp)
		if err == nil && errResp.Error != nil {
			err = errors.New("Error: " + errResp.Error.Type + " - " + errResp.Error.Message)
		}
	}
	return jsonObj, err
}
