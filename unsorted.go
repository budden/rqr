package main

import (
	"net/http"
	"strings"
)

func return404IfExtraURLChars(path string, w http.ResponseWriter, req *http.Request) (doReturn bool) {
	if strings.TrimPrefix(req.URL.Path, path) != "" {
		w.WriteHeader(http.StatusNotFound)
		doReturn = true
	}
	return
}
