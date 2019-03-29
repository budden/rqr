package main

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/budden/rqr/pkg/errorcodes"
)

func executeFetchTask(pt *ParsedFetchTask) (et *ExecutedFetchTask, err *errorWithCode) {
	var b bytes.Buffer
	request, err1 := http.NewRequest(pt.Method, pt.URL, &b)
	if err1 != nil {
		err = newErrorWithCode(errorcodes.FailedToMakeARequest, "%v", err1)
		return
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
	et = &ExecutedFetchTask{Httpstatus: resp.StatusCode, Headers: resp.Header, Bodylength: len(body)}
	return
}
