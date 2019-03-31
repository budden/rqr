package main

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/budden/rqr/pkg/errorcodes"
	"golang.org/x/net/netutil"
)

const (
	fetchTaskGetURL    = "/fetchtaskget/"
	fetchTaskDeleteURL = "/fetchtaskdelete/"
)

func main() {
	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/fetchtaskadd", handleFetchTaskAdd)
	http.HandleFunc("/fetchtasklist", handleFetchTaskList)
	http.HandleFunc(fetchTaskGetURL, handleFetchTaskGet)
	http.HandleFunc(fetchTaskDeleteURL, handleFetchTaskDelete)

	// https://habr.com/ru/post/197468/
	ThisHTTPServer := &http.Server{
		Addr:           ":8086",
		Handler:        http.DefaultServeMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20}

	listener, err := net.Listen("tcp", ThisHTTPServer.Addr)

	if err != nil {
		log.Fatalln(err)
	}

	limitListener := netutil.LimitListener(listener, 10)

	log.Fatalln(ThisHTTPServer.Serve(limitListener))

}

func handleRoot(w http.ResponseWriter, req *http.Request) {
	if checkNoExtraURLChars("/", w, req) || checkHTTPMethod("GET", w, req) {
		return
	}
	WriteReplyToResponseAsJSON(w, req, errorcodes.OK, []string{
		"Requester service.",
		"Use POST /fetchtaskadd json urlencoded to add a fetch task",
		"Use GET /fetchtaskget/ID to get a fetch task",
		"Use GET /fetchtasklist?offset=N&limit=N to get a list (both params are optional)",
		"Use POST /fetchtaskdelete/ID to delete a fetch task",
		"Use GET / to obtain this help",
		"Replies are always with Content-type = application/json"})
	return
}
