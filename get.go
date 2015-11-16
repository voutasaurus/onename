package main

import (
	"errors"
	"io/ioutil"
	"net/http"
)

// GetRequest GETs the json response using client c
func (c *client) GetRequest(url string) ([]byte, error) {
	respJSON := []byte{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return respJSON, errors.New("Request not sent. Could not construct request.\n\t Check the URL \"" + url + "\" - got error: " + err.Error())
	}

	req.SetBasicAuth(c.ApiID, c.ApiSecret)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return respJSON, errors.New("Request failed.\n\t Got error: " + err.Error())
	}

	defer resp.Body.Close()

	respJSON, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return respJSON, errors.New("Request completed but read of response body failed.\n\t Got error: " + err.Error())
	}

	return respJSON, err
}
