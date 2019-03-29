package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

// fetchtaskget, fetchtasklist

func handleFetchTaskList(w http.ResponseWriter, req *http.Request) {
	if return500IfNotMethod("GET", w, req) {
		return
	}
	unsorted := allFetchTasks()
	sorted := unsorted[:]
	// sort them (and destroy unsorted)
	sort.Slice(sorted, func(i, j int) bool {
		return FetchTaskLessThan(sorted[i], sorted[j])
	})
	err := req.ParseForm()
	if reportFetchTaskErrorToClientIf(err, w) {
		return
	}
	offset, _, doReturn1 := GetZeroOrOneNonNegativeIntFormValueOrReportAnError("offset", w, req)
	if doReturn1 {
		return
	}
	limit, _, doReturn2 := GetZeroOrOneNonNegativeIntFormValueOrReportAnError("limit", w, req)
	if doReturn2 {
		return
	}
	length := len(sorted)
	beg, end := startAndLimitToBegAndEnd(offset, limit, length)
	selected := sorted[beg:end]
	results := make([]*FetchTaskAsJSON, end-beg)
	for i, task := range selected {
		results[i] = convertFetchTaskToJSON(task)
	}
	encoder := json.NewEncoder(w)
	err = encoder.Encode(results)
	if err != nil {
		// FIXME: It is too late to set HTTP status. Should we serialize to a string and then
		// write a string?
		log.Printf("Failed to encode results, error is %#v", err)
	}
}

// for paging
func startAndLimitToBegAndEnd(start, limit, length int) (beg, end int) {
	beg = start
	if beg > length {
		beg = length
	}
	if limit == 0 {
		limit = length
	}
	end = beg + limit
	if end > length {
		end = length
	}
	return
}

func handleFetchTaskGet(w http.ResponseWriter, req *http.Request) {
	if return500IfNotMethod("GET", w, req) {
		return
	}
	_, ft, doReturn := getFetchTaskFromLastURLSegment(fetchTaskGetURL, w, req)
	if doReturn {
		return
	}
	encoder := json.NewEncoder(w)
	ftj := convertFetchTaskToJSON(ft)
	err := encoder.Encode(ftj)
	if err != nil {
		// FIXME: It is too late to set HTTP status. Should we serialize to a string and then
		// write a string?
		log.Printf("Failed to encode results, error is %#v", err)
	}
}
