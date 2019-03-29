package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"

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

// https://golang.org/doc/faq#nil_error recommends to always return error to avoid confusion
// We do prefer to return typed error instead of error because it adds more static typing.
// Maybe it is wrong in terms of performance... For that we could create a sort of
// isZeroError<T> for types we want.
func isZeroError(err error) bool {
	return (err == nil || reflect.ValueOf(err) == reflect.Zero(reflect.TypeOf(err)))
}

func reportFetchTaskErrorToClientIf(err error, w http.ResponseWriter) (doReturn bool) {
	if isZeroError(err) {
		return
	}
	doReturn = true
	w.WriteHeader(http.StatusBadRequest)
	// here "&& je != nil" handles nil_error case
	if je, ok := err.(*errorWithCode); ok && je != nil {
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
