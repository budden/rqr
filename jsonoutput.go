package main

// When the architecture is not perfect (yet), we put homeless things here

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

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
