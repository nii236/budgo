package main

import (
	"encoding/json"
	"fmt"
)

// Error is an error type for easy JSON marshalling
type Error struct {
	Message string `json:"message,omitempty"`
}

// Err is the constructor for the Error struct
func Err(err error) *Error {
	return &Error{
		Message: err.Error(),
	}
}

func (e *Error) Error() string {
	return e.Message
}

// JSON returns the JSON representation of the error struct
func (e *Error) JSON() string {
	b, err := json.Marshal(e)
	if err != nil {
		return "could not marshal error struct: " + err.Error()
	}
	fmt.Println(e.Message)
	return string(b)
}
