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

type OneNameResponse struct {
	Error *OneNameError `json:"error,omitempty"`
}

type OneNameError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// GET
func (c *client) GetRequest(url string) (OneNameResponse, error) {
	responseObject := OneNameResponse{}

	req, err := http.NewRequest("GET", url, nil)
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

func (c *client) GetUsers(usernames []string) (OneNameResponse, error) {
	url := c.BaseUrl + "/users/" + strings.Join(usernames, ",")
	return c.GetRequest(url)
}

func (c *client) SearchUsers(query string) (OneNameResponse, error) {
	url := c.BaseUrl + "/search?query=" + query
	return c.GetRequest(url)
}

func (c *client) GetAllUsers() (OneNameResponse, error) {
	url := c.BaseUrl + "/users"
	return c.GetRequest(url)
}

func (c *client) GetUserStats() (OneNameResponse, error) {
	url := c.BaseUrl + "/stats/users"
	return c.GetRequest(url)
}

func (c *client) GetUnspents(address string) (OneNameResponse, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return OneNameResponse{}, errors.New("Address must be a valid cryptocurrency address. Instead it is " + address + err.Error())
	}
	url := c.BaseUrl + "/addresses/" + address + "/unspents"
	return c.GetRequest(url)
}

func (c *client) GetNames(address string) (OneNameResponse, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return OneNameResponse{}, errors.New("Address must be a valid cryptocurrency address. Instead it is " + address + err.Error())
	}
	url := c.BaseUrl + "/addresses/" + address + "/names"
	return c.GetRequest(url)
}

func (c *client) GetDKIMInfo(domain string) (OneNameResponse, error) {
	url := c.BaseUrl + "/domains/" + domain + "/dkim"
	return c.GetRequest(url)
}

// POST
type OneNamePost struct {
	Passname         string `json:"passname,omitempty"`
	RecipientAddress string `json:"recipient_address,omitempty"`
	SignedHex        string `json:"signed_hex,omitempty"`
}

func (c *client) PostRequest(url string, payload OneNamePost) (OneNameResponse, error) {
	responseObject := OneNameResponse{}

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

func (c *client) RegisterUser(username, address string) (OneNameResponse, error) {
	_, _, err := base58.CheckDecode(address)
	if err != nil {
		return OneNameResponse{}, err
	}
	url := c.BaseUrl + "/users"
	payload := OneNamePost{username, address, ""}
	return c.PostRequest(url, payload)
}

func (c *client) BroadcastTransactions(signedHex string) (OneNameResponse, error) {
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
	resp, err := c.GetRequest("https://api.onename.com/v1/users")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)

}
