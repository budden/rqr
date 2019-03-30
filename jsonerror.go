package main

import (
	"fmt"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

// ErrorWithCodeData is an error which is returned to the client in JSON format
type ErrorWithCodeData struct {
	ECode  errorcodes.FetchTaskErrorCode
	EError string
}

// ErrorWithCode is for an experimental typed error handling. We define an
// interface which is a superset of `error` and use it at some places instead
// of error, to get richer error objects. Let us note that
// https://golang.org/doc/faq#nil_error recommends to always return error
// to avoid "nil interface" issue. But always returning error is not very
// statically typed, so we're trying to make a step towards more declarative
// function signatures which declare the type of possible errors.
type ErrorWithCode interface {
	Code() errorcodes.FetchTaskErrorCode
	Error() string
	Data() *ErrorWithCodeData
}

// Code ...
func (ewc *ErrorWithCodeData) Code() errorcodes.FetchTaskErrorCode {
	return ewc.ECode
}

func (ewc *ErrorWithCodeData) Error() string {
	return ewc.EError
}

// Data ...
func (ewc *ErrorWithCodeData) Data() *ErrorWithCodeData {
	return ewc
}

func newErrorWithCode(ECode errorcodes.FetchTaskErrorCode, format string, args ...interface{}) ErrorWithCode {
	return &ErrorWithCodeData{ECode: ECode, EError: fmt.Sprintf(format, args...)}
}

// jsonFetchTask type expresses the fact that fetchTask must be a JSON array
type jsonFetchTask []string

func reportFetchTaskErrorToClientIf(err error, w http.ResponseWriter, req *http.Request) (
	doReturn bool) {
	if err == nil {
		return
	}
	doReturn = true
	status := errorcodes.UnknownError
	if je, ok := err.(ErrorWithCode); ok {
		status = je.Code()
	}
	// doReturn is alreay true, no need to assign
	WriteReplyToResponseAsJSON(w, req, status, err.Error())
	return
}
