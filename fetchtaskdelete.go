package main

import (
	"net/http"
)

func handleFetchTaskDelete(w http.ResponseWriter, req *http.Request) {
	if return500IfNotMethod("POST", w, req) {
		return
	}
	ID, _, doReturn := getFetchTaskFromLastURLSegment(fetchTaskDeleteURL, w, req)
	if doReturn {
		return
	}
	err := eraseFetchTask(ID)
	if reportFetchTaskErrorToClientIf(err, w) {
		return
	}
	w.WriteHeader(http.StatusOK)
}
