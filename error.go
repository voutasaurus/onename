package main

import "errors"

type OneNameErrorResponse struct {
	Error *OneNameError `json:"error,omitempty"`
}

type OneNameError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (e OneNameError) String() string {
	return "Error: " + e.Type + " - " + e.Message
}

func JSONRead(err error, jsonBytes []byte) error {
	if len(jsonBytes) > 0 {
		return errors.New("Request+Response complete but could not read response as JSON.\n\t Response was: " + string(jsonBytes) + "\t Got error: " + err.Error())
	}
	return errors.New("Request+Response complete but could not read response as JSON.\n\t Response was empty. Got error: " + err.Error())
}
