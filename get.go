package main

import (
	"io/ioutil"
	"net/http"
)

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
