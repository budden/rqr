package main

// When the architecture is not perfect (yet), we put homeless things here

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/budden/rqr/pkg/errorcodes"
	"github.com/pkg/errors"
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

func failIfMethodIsNot(method string, w http.ResponseWriter, req *http.Request) (doReturn bool) {
	if req.Method != method {
		WriteReplyToResponseAsJSON(w, req, errorcodes.IncorrectRequestMethod, nil)
		doReturn = true
	}
	return
}

// SetJSONContentType sets Content-Type application/json
func SetJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// WriteReplyToResponseAsJSON handles all errors
func WriteReplyToResponseAsJSON(
	w http.ResponseWriter,
	req *http.Request,
	status errorcodes.FetchTaskErrorCode,
	contents interface{}) (doReturn bool) {
	dataToEncode := JSONTopLevel{Status: status, Statustext: status.String(), Contents: contents}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(dataToEncode)
	if err != nil {
		// we don't put the data to the log. Data can be huge and it can be broken. Instead, we put a
		// truncated URL
		url := req.URL.Path
		if len(url) > 256 {
			url = url[0:(256-3)] + "..."
		}
		log.Printf("Failed to encode results, URL is «%v», status is %v, error is %#v",
			url, status.String(), err)
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
		err := newErrorWithCode(errorcodes.DuplicateParamInTheForm, "Query parameter «%s» is duplicated", key)
		reportFetchTaskErrorToClientIf(err, w, req)
		doReturn = true
		return
	}
	valueS := values[0]
	if valueS == "" {
		return
	}
	var err error
	value, err = strconv.Atoi(valueS)
	if reportFetchTaskErrorToClientIf(err, w, req) {
		doReturn = true
		return
	}
	if value < 0 {
		err := fmt.Errorf("Negative value of parameter «%s» is not allowed", key)
		reportFetchTaskErrorToClientIf(err, w, req)
		doReturn = true
		return
	}
	return
}

// this one is shared between fetchtaskget and fetchtaskdelete, so it is
// handy to pack into the function
func getFetchTaskFromLastURLSegment(URLBase string, w http.ResponseWriter, req *http.Request) (
	ID string, ft *FetchTask, doReturn bool) {
	doReturn = true
	ID = strings.TrimPrefix(req.URL.Path, URLBase)
	matched, err := regexp.Match("^[0-9]+$", []byte(ID))
	if err != nil {
		err = errors.Wrapf(err, "Error when matching regexp")
		reportFetchTaskErrorToClientIf(err, w, req)
		ID = ""
		return
	}
	if !matched {
		err = newErrorWithCode(errorcodes.IncorrectIDFormat, "")
		reportFetchTaskErrorToClientIf(err, w, req)
		ID = ""
		return
	}
	ft, ok := getFetchTask(ID)
	if !ok {
		ID = ""
		err = newErrorWithCode(errorcodes.FetchTaskNotFound, "")
		reportFetchTaskErrorToClientIf(err, w, req)
		return
	}
	doReturn = false
	return
}
