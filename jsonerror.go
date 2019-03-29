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
	Code    errorcodes.FetchTaskErrorCode
	Message string
}

func newErrorWithCode(Code errorcodes.FetchTaskErrorCode, format string, args ...interface{}) *errorWithCode {
	return &errorWithCode{Code: Code, Message: fmt.Sprintf(format, args...)}
}

// Error ...
func (je *errorWithCode) Error() string {
	return je.Message
}

// jsonFetchTask type expresses the fact that fetchTask must be a JSON array
type jsonFetchTask []string

func reportFetchTaskErrorToClientIf(err error, w http.ResponseWriter) (doReturn bool) {
	if err == nil {
		return
	}
	doReturn = true
	w.WriteHeader(http.StatusBadRequest)
	if je, ok := err.(*errorWithCode); ok {
		// Let's add a textual representation of a error code.
		errorAndStringCode := []interface{}{je.Code.String(), je}
		encoder := json.NewEncoder(w)
		err := encoder.Encode(errorAndStringCode)
		if err != nil {
			log.Printf("Error while sending error response to a client: %#v\n", err)
		}
	}
	return
}
