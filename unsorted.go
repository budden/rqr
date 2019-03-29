package main

// When the architecture is not perfect (yet), we put homeless things here

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/budden/rqr/pkg/errorcodes"
)

func return404IfExtraURLChars(path string, w http.ResponseWriter, req *http.Request) (doReturn bool) {
	if strings.TrimPrefix(req.URL.Path, path) != "" {
		w.WriteHeader(http.StatusNotFound)
		doReturn = true
	}
	return
}

func return500IfNotMethod(method string, w http.ResponseWriter, req *http.Request) (doReturn bool) {
	if req.Method != method {
		w.WriteHeader(http.StatusBadRequest)
		doReturn = true
	}
	return
}

// GetZeroOrOneNonNegativeIntFormValueOrReportAnError extracts an integer value from the req.Form.
// req.ParseForm() must be called before this one. If there are none, ok == nil and value = 0
// If there are more than one, responses with http.StatusBadRequest. In case of error, sets doReturn to true.
func GetZeroOrOneNonNegativeIntFormValueOrReportAnError(key string, w http.ResponseWriter, req *http.Request) (
	value int, ok, doReturn bool) {
	values, ok1 := req.Form[key]
	if !ok1 || len(values) == 0 {
		return
	}
	if len(values) > 1 {
		err := newErrorWithCode(errorcodes.DuplicatedValueInTheForm, "Query parameter «%s» is duplicated", key)
		reportFetchTaskErrorToClientIf(err, w)
		doReturn = true
		return
	}
	valueS := values[0]
	if valueS == "" {
		return
	}
	var err error
	value, err = strconv.Atoi(valueS)
	if reportFetchTaskErrorToClientIf(err, w) {
		doReturn = true
		return
	}
	if value < 0 {
		err := fmt.Errorf("Negative value of parameter «%s» is not allowed", key)
		reportFetchTaskErrorToClientIf(err, w)
		doReturn = true
		return
	}
	return
}
