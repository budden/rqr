package main

import (
	"net/http"
	"sort"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleFetchTaskList(w http.ResponseWriter, req *http.Request) {
	SetJSONContentType(w)
	if failIfMethodIsNot("GET", w, req) {
		return
	}
	unsorted := allFetchTasks()
	sorted := unsorted[:]
	// sort them (and destroy unsorted)
	sort.Slice(sorted, func(i, j int) bool {
		return FetchTaskLessThan(sorted[i], sorted[j])
	})
	err := req.ParseForm()
	if reportFetchTaskErrorToClientIf(err, w, req) {
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
	records := make([]*FetchTaskAsJSON, end-beg)
	for i, task := range selected {
		records[i] = convertFetchTaskToJSON(task)
	}
	result := &FetchTaskListAsJSON{Length: length, Records: records}
	// no need to analyze wasError, we're exiting anyways
	_ = WriteReplyToResponseAsJSON(w, req, errorcodes.OK, result)
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
