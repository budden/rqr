package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/budden/rqr/pkg/errorcodes"
)

func executeFetchTask(pt *ParsedFetchTask) (et *ExecutedFetchTask, err ErrorWithCode) {
	var r io.Reader
	if pt.Body != "" {
		// avoid buffer allocation, https://habr.com/ru/post/307554/
		r = strings.NewReader(pt.Body)
	}
	request, err1 := http.NewRequest(pt.Method, pt.URL, r)
	if err1 != nil {
		err = newErrorWithCode(errorcodes.FailedToMakeARequest, "%v", err1)
		return
	}
	for k, v := range pt.Headers {
		// https://stackoverflow.com/a/41034588/9469533
		if k == "Host" {
			request.Host = v
		} else {
			request.Header.Set(k, v)
		}
	}
	resp, err2 := http.DefaultClient.Do(request)
	if err2 != nil {
		err = newErrorWithCode(errorcodes.FailedToSendARequest, "%v", err2)
		return
	}
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		err = newErrorWithCode(errorcodes.FailedToReadRequestBody, "%v", err3)
		return
	}
	// For debugging purposes we can print a body here
	// fmt.Printf("Body: «%s»\n", string(body))
	et = &ExecutedFetchTask{Httpstatus: resp.StatusCode, Headers: resp.Header, Bodylength: len(body)}
	return
}
