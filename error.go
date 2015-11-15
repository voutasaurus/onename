package main

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
