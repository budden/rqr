package main

import (
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleFetchTaskDelete(w http.ResponseWriter, req *http.Request) {
	SetJSONContentType(w)
	if checkHTTPMethod("POST", w, req) {
		return
	}
	ID, _, wasError := getFetchTaskFromLastURLSegment(fetchTaskDeleteURL, w, req)
	if wasError {
		return
	}
	err := eraseFetchTask(ID)
	if reportFetchTaskErrorToClientIf(err, w, req) {
		return
	}
	_ = WriteReplyToResponseAsJSON(w, req, errorcodes.OK, nil)
}
