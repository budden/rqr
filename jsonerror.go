package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

// jsonError is an error which is returned to the client in JSON format
type jsonError struct {
	Code    errorcodes.TaskErrorCode
	Message string
}

// Error ...
func (je *jsonError) Error() string {
	return je.Message
}

// jsonTask type expresses the fact that task must be a JSON array
type jsonTask []string

func reportTaskErrorToClientIf(err error, w http.ResponseWriter) (doReturn bool) {
	if err == nil {
		return
	}
	doReturn = true
	if je, ok := err.(*jsonError); ok {
		encoder := json.NewEncoder(w)
		err := encoder.Encode(je)
		if err != nil {
			log.Printf("Error while sending error response to a client: %#v\n", err)
		}
	}
	return
}
