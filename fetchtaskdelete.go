package main

import (
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleFetchTaskDelete(w http.ResponseWriter, req *http.Request) {
	SetJSONContentType(w)
	if failIfMethodIsNot("POST", w, req) {
		return
	}
	ID, _, doReturn := getFetchTaskFromLastURLSegment(fetchTaskDeleteURL, w, req)
	if doReturn {
		return
	}
	err := eraseFetchTask(ID)
	if reportFetchTaskErrorToClientIf(err, w, req) {
		return
	}
	_ = WriteReplyToResponseAsJSON(w, req, errorcodes.OK, nil)
}
