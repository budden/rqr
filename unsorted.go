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

func checkNoExtraURLChars(path string, w http.ResponseWriter, req *http.Request) (wasError bool) {
	if strings.TrimPrefix(req.URL.Path, path) != "" {
		WriteReplyToResponseAsJSON(w, req, errorcodes.IncorrectURL, "POST to / to obtain a help on correct URLs")
		wasError = true
	}
	return
}

func checkHTTPMethod(method string, w http.ResponseWriter, req *http.Request) (wasError bool) {
	if req.Method != method {
		WriteReplyToResponseAsJSON(w, req, errorcodes.IncorrectRequestMethod, "POST to / to obtain a help on correct URLs")
		wasError = true
	}
	return
}

func failIfMethodIsNot(method string, w http.ResponseWriter, req *http.Request) (wasError bool) {
	if req.Method != method {
		WriteReplyToResponseAsJSON(w, req, errorcodes.IncorrectRequestMethod, nil)
		wasError = true
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
	contents interface{}) (wasError bool) {
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
		wasError = true
	}
	return
}

// GetZeroOrOneNonNegativeIntFormValueOrReportAnError extracts an integer value from the req.Form.
// req.ParseForm() must be called before this one. If there are none, ok == nil and value = 0
// If there are more than one, responses with http.StatusBadRequest. In case of error, sets wasError to true.
func GetZeroOrOneNonNegativeIntFormValueOrReportAnError(key string, w http.ResponseWriter, req *http.Request) (
	value int, ok, wasError bool) {
	values, ok1 := req.Form[key]
	if !ok1 || len(values) == 0 {
		return
	}
	if len(values) > 1 {
		err := newErrorWithCode(errorcodes.DuplicateParamInTheForm, "Query parameter «%s» is duplicated", key)
		reportFetchTaskErrorToClientIf(err, w, req)
		wasError = true
		return
	}
	valueS := values[0]
	if valueS == "" {
		return
	}
	var err error
	value, err = strconv.Atoi(valueS)
	if reportFetchTaskErrorToClientIf(err, w, req) {
		wasError = true
		return
	}
	if value < 0 {
		err := fmt.Errorf("Negative value of parameter «%s» is not allowed", key)
		reportFetchTaskErrorToClientIf(err, w, req)
		wasError = true
		return
	}
	return
}

// this one is shared between fetchtaskget and fetchtaskdelete, so it is
// handy to pack into the function
func getFetchTaskFromLastURLSegment(URLBase string, w http.ResponseWriter, req *http.Request) (
	ID string, ft *FetchTask, wasError bool) {
	wasError = true
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
	wasError = false
	return
}
