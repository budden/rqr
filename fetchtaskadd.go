package main

import (
	"encoding/json"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

// https://stackoverflow.com/a/15685432/9469533
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
	var JSONinput interface{}

	failParseIf := func(condition bool, format string, args ...interface{}) (wasError bool) {
		if condition {
			err = newErrorWithCode(errorcodes.FailedToParsefetchTaskJSON, format, args...)
			wasError = true
		}
		return
	}

	err = decoder.Decode(&JSONinput)
	// this is not an efficient way to check errors, but it saves lines of code :)
	if failParseIf(err != nil, "Failed to parse request JSON data. Error is %#v", err) {
		return
	}
	JSONarray, ok1 := JSONinput.([]interface{})
	if failParseIf(!ok1, "fetch task data must be an array") {
		return
	}
	lenFetchTask := len(JSONarray)
	if failParseIf(lenFetchTask != 2 && lenFetchTask != 4,
		"JSON fetchTask must be of the form [method, address] or of the form [method, address, headers, body]") {
		return
	}
	method, ok2 := JSONarray[0].(string)
	if failParseIf(!ok2, "first element of JSON fetch task array must be a string (method)") {
		return
	}
	url, ok3 := JSONarray[1].(string)
	if failParseIf(!ok3, "second element of JSON fetch task array must be a string (URL)") {
		return
	}
	result := &ParsedFetchTask{Method: method, URL: url}
	if lenFetchTask == 4 {
		headers, ok4 := JSONarray[2].(map[string]interface{})
		if failParseIf(!ok4, "third element of JSON fetch task array, if present, must be an object with string fields (header)") {
			return
		}
		headersStrings := map[string]string{}
		for k, v := range headers {
			vString, ok5 := v.(string)
			if failParseIf(!ok5, "If headers present, they must be an object with string fields") {
				return
			}
			headersStrings[k] = vString
		}
		body, ok6 := JSONarray[3].(string)
		if failParseIf(!ok6, "fourth element of JSON fetch task array, if present, must be a string (body)") {
			return
		}
		result.Headers = headersStrings
		result.Body = body
	}
	// all previous return statements were due to errors, so we take care that
	// an empty result is returned in case of error. Now we
	// are safe to return actual data.
	pt = result
	return
}
