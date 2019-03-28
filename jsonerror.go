package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

// errorWithCode is an error which is returned to the client in JSON format
type errorWithCode struct {
	Code    errorcodes.TaskErrorCode
	Message string
}

func newErrorWithCode(Code errorcodes.TaskErrorCode, format string, args ...interface{}) *errorWithCode {
	return &errorWithCode{Code: Code, Message: fmt.Sprintf(format, args...)}
}

// Error ...
func (je *errorWithCode) Error() string {
	return je.Message
}

// jsonTask type expresses the fact that task must be a JSON array
type jsonTask []string

func reportTaskErrorToClientIf(err error, w http.ResponseWriter) (doReturn bool) {
	if err == nil {
		return
	}
	doReturn = true
	if je, ok := err.(*errorWithCode); ok {
		encoder := json.NewEncoder(w)
		err := encoder.Encode(je)
		if err != nil {
			log.Printf("Error while sending error response to a client: %#v\n", err)
		}
	}
	return
}
