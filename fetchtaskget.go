package main

import (
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleFetchTaskGet(w http.ResponseWriter, req *http.Request) {
	SetJSONContentType(w)
	if failIfMethodIsNot("GET", w, req) {
		return
	}
	_, ft, doReturn := getFetchTaskFromLastURLSegment(fetchTaskGetURL, w, req)
	if doReturn {
		return
	}
	ftj := convertFetchTaskToJSON(ft)
	_ = WriteReplyToResponseAsJSON(w, req, errorcodes.OK, ftj)
}
