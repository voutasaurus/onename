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

	"github.com/btcsuite/btcutil/base58"
)

var BASE_URL = "https://api.onename.com/v1"

type client struct {
	ApiID     string
	ApiSecret string
	BaseUrl   string
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

	users, err := c.GetUsers([]string{"fredwilson"})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(users)
	}

	stats, err := c.GetUserStats()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(stats)
	}

	searchResults, err := c.SearchUsers("wenger")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(searchResults)

	unspentsResults, err := c.GetUnspents("1QHDGGLEKK7FZWsBEL78acV9edGCTarqXt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(unspentsResults)

	namesResults, err := c.GetNames("1QHDGGLEKK7FZWsBEL78acV9edGCTarqXt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(namesResults)

	dkimResult, err := c.GetDKIMInfo("onename.com")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(dkimResult)

	return

}
