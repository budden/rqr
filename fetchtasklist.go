package main

import (
	"net/http"
	"sort"

	"github.com/budden/rqr/pkg/errorcodes"
)

func handleFetchTaskList(w http.ResponseWriter, req *http.Request) {
	SetJSONContentType(w)
	if checkHTTPMethod("GET", w, req) {
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

// for paging. Start and Limit are assumed to be non-negative w/o checks. In our requests,
// negative values are rejected by the GetZeroOrOneNonNegativeIntFormValueOrReportAnError
// Zero limit means no limit.
func startAndLimitToBegAndEnd(nonNegativeStart, nonNegativeLimit, length int) (beg, end int) {
	beg = nonNegativeStart
	if beg > length {
		beg = length
	}
	if nonNegativeLimit == 0 {
		nonNegativeLimit = length
	}
	end = beg + nonNegativeLimit
	if end > length {
		end = length
	}
	return
}
