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

// WriteReplyToResponseAsJSON writes it's errors to log (obviously we can't sent them to client)
// and returns wasError if there was an error
func WriteReplyToResponseAsJSON(
	w http.ResponseWriter, req *http.Request,
	status errorcodes.FetchTaskErrorCode, contents interface{}) (wasError bool) {
	dataToEncode := JSONTopLevel{Status: status, Statustext: status.String(), Contents: contents}
	encoder := json.NewEncoder(w)
	err := encoder.Encode(dataToEncode)
	if err != nil {
		// we don't put the data to the log. Data can be huge, broken and contain confidential info.
		// Instead, we put a truncated URL (URL can also be huge)
		url := req.URL.Path
		url = trimToTheNumberOfRunes(url, 256)
		log.Printf("Failed to encode results, URL is «%v», status is %v, error is %#v",
			url, status.String(), err)
		wasError = true
	}
	return
}
