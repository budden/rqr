package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

// https://stackoverflow.com/a/15685432/9469533
// To test, use curl -i -X POST -d "[\"GET\", \"google.com\"]" http://localhost:8086/fetchTaskadd
// To test error reporting, remove the comma from JSON :)
func handleFetchTaskAdd(w http.ResponseWriter, req *http.Request) {
	if return500IfNotMethod("POST", w, req) {
		return
	}
	pt, err := convertJSONFetchTaskToParsedFetchTask(req)
	if reportFetchTaskErrorToClientIf(err, w) {
		return
	}
	et, err1 := executeFetchTask(pt)
	if reportFetchTaskErrorToClientIf(err1, w) {
		return
	}
	ft := saveFetchTask(pt, et)          // no expected errors here
	ftJSON := convertFetchTaskToJSON(ft) // no expected errors here
	sendFetchTask(w, ftJSON)             // errors are handled inside
}

func sendFetchTask(w http.ResponseWriter, ftJSON *FetchTaskAsJSON) {
	// we call w.WriteHeader everywhere to get a warning in case of multiple resonse.WriteHeader calls.
	// unfortunately there is no error to check.
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(ftJSON)
	if err != nil {
		// to late to change status
		log.Printf("Failed to encode response to json, respons is %+v, error is %v", ftJSON, err)
	}
}

func convertJSONFetchTaskToParsedFetchTask(req *http.Request) (pt *ParsedFetchTask, err error) {
	decoder := json.NewDecoder(req.Body)
	ji := jsonFetchTask{}
	err = decoder.Decode(&ji)
	// this is not an efficient way to check errors, but it saves lines of code :)

	if err != nil {
		err = newErrorWithCode(errorcodes.FailedToParsefetchTaskJSON, "Failed to parse request JSON data. Error is %#v", err)
		return
	}
	lenFetchTask := len(ji)
	if lenFetchTask != 2 && lenFetchTask != 4 {
		err = newErrorWithCode(errorcodes.FailedToParsefetchTaskJSON,
			"JSON fetchTask must be of the form [method, address] or of the form [method, address, headers, body]")
		return
	}
	pt = &ParsedFetchTask{Method: ji[0], URL: ji[1]}
	if lenFetchTask == 4 {
		pt.Headers = ji[2]
		pt.Body = ji[3]
	}
	return
}