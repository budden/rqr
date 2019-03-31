package main

import (
	"encoding/json"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

// https://stackoverflow.com/a/15685432/9469533
// To test, use curl -i -X POST -d "[\"GET\", \"google.com\"]" http://localhost:8086/fetchTaskadd
// To test error reporting, remove the comma from JSON :)
func handleFetchTaskAdd(w http.ResponseWriter, req *http.Request) {
	SetJSONContentType(w)
	if checkHTTPMethod("POST", w, req) {
		return
	}
	pt, err := convertJSONFetchTaskToParsedFetchTask(req)
	if reportFetchTaskErrorToClientIf(err, w, req) {
		return
	}
	et, err1 := executeFetchTask(pt)
	if reportFetchTaskErrorToClientIf(err1, w, req) {
		return
	}
	ft := saveFetchTask(pt, et)          // no expected errors here
	ftJSON := convertFetchTaskToJSON(ft) // no expected errors here
	// no need to check failure here, we're exiting anyways
	_ = WriteReplyToResponseAsJSON(w, req, errorcodes.OK, ftJSON)
}

func convertJSONFetchTaskToParsedFetchTask(req *http.Request) (pt *ParsedFetchTask, err error) {
	decoder := json.NewDecoder(req.Body)
	var ji interface{}

	failParseIf := func(condition bool, format string, args ...interface{}) (wasError bool) {
		if condition {
			err = newErrorWithCode(errorcodes.FailedToParsefetchTaskJSON, format, args...)
			wasError = true
		}
		return
	}

	err = decoder.Decode(&ji)
	// this is not an efficient way to check errors, but it saves lines of code :)

	if failParseIf(err != nil, "Failed to parse request JSON data. Error is %#v", err) {
		return
	}
	ja, ok1 := ji.([]interface{})
	if failParseIf(!ok1, "fetch task data must be an array") {
		return
	}
	lenFetchTask := len(ja)
	if failParseIf(lenFetchTask != 2 && lenFetchTask != 4,
		"JSON fetchTask must be of the form [method, address] or of the form [method, address, headers, body]") {
		return
	}
	method, ok2 := ja[0].(string)
	if failParseIf(!ok2, "first element of JSON fetch task array must be a string (method)") {
		return
	}
	url, ok3 := ja[1].(string)
	if failParseIf(!ok3, "second element of JSON fetch task array must be a string (URL)") {
		return
	}
	pt = &ParsedFetchTask{Method: method, URL: url}
	if lenFetchTask == 4 {
		headers, ok4 := ja[2].(map[string]interface{})
		if failParseIf(!ok4, "third element of JSON fetch task array, if present, must be an object with string fields (header)") {
			return
		}
		headersStrings := map[string]string{}
		for k, v := range headers {
			vString, ok5 := v.(string)
			if failParseIf(!ok5, "If headers present, it must be an object with string fields") {
				return
			}
			headersStrings[k] = vString
		}
		body, ok6 := ja[3].(string)
		if failParseIf(!ok6, "fourth element of JSON fetch task array, if present, must be a string (body)") {
			return
		}
		pt.Headers = headersStrings
		pt.Body = body
	}
	return
}
