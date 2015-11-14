package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

var BASE_URL = "https://api.onename.com/v1"

type client struct {
	ApiID     string
	ApiSecret string
	BaseUrl   string
}

type OneNameErrorResponse struct {
	Error *OneNameError `json:"error,omitempty"`
}

type OneNameError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// GetRequest GETs the json response using client c
func (c *client) GetRequest(url string) ([]byte, error) {
	respJSON := []byte{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return respJSON, err
	}

	req.SetBasicAuth(c.ApiID, c.ApiSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return respJSON, err
	}

	defer resp.Body.Close()

	respJSON, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return respJSON, err
	}

	return respJSON, err
}

// GetUserObject parses the JSON responses matching the Users type
func (c *client) GetUserObjects(url string) (Users, error) {
	var jsonObj Users
	jsonBytes, err := c.GetRequest(url)
	if err != nil {
		return jsonObj, err
	}
	err = json.Unmarshal(jsonBytes, &jsonObj)
	if err != nil {
		return jsonObj, err
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

func (c *client) GetUsers(usernames []string) (Users, error) {
	url := c.BaseUrl + "/users/" + strings.Join(usernames, ",")
	return c.GetUserObjects(url)
}

func (c *client) SearchUsers(query string) (Users, error) {
	url := c.BaseUrl + "/search?query=" + query
	return c.GetUserObjects(url)
}

func (c *client) GetAllUsers() (Users, error) {
	url := c.BaseUrl + "/users"
	return c.GetUserObjects(url)
}

func (c *client) GetUserStats() ([]byte, error) {
	url := c.BaseUrl + "/stats/users"
	return c.GetRequest(url)
}

func (c *client) GetUnspents(address string) ([]byte, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return []byte{}, errors.New("Address must be a valid cryptocurrency address. Instead it is " + address + err.Error())
	}
	url := c.BaseUrl + "/addresses/" + address + "/unspents"
	return c.GetRequest(url)
}

func (c *client) GetNames(address string) ([]byte, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return []byte{}, errors.New("Address must be a valid cryptocurrency address. Instead it is " + address + err.Error())
	}
	url := c.BaseUrl + "/addresses/" + address + "/names"
	return c.GetRequest(url)
}

func (c *client) GetDKIMInfo(domain string) ([]byte, error) {
	url := c.BaseUrl + "/domains/" + domain + "/dkim"
	return c.GetRequest(url)
}

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

func main() {

	file, err := os.Open("cred.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	c := client{lines[0], lines[1], "https://api.onename.com/v1"}
	/*	resp, err := c.GetRequest("https://api.onename.com/v1/users")
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)*/

	users, err := c.GetAllUsers()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(users)
	return

}
