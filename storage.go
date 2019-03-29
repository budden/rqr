package main

import (
	"math/big"

	"github.com/budden/rqr/pkg/errorcodes"
)

var queryID = big.NewInt(0)
var one = big.NewInt(1)

var fetchTaskStorage map[string]*FetchTask

func saveFetchTask(pt *ParsedFetchTask, et *ExecutedFetchTask) (t *FetchTask) {
	// FIXME will the map be of size 1?
	queryID.Add(queryID, one)
	iString := queryID.String()
	fetchTaskStorage[iString] = &FetchTask{ID: iString, pt: pt, et: et}
	return
}

func getFetchTask(iString string) (t *FetchTask, ok bool) {
	t, ok = fetchTaskStorage[iString]
	return
}

func eraseFetchTask(iString string) (err *errorWithCode) {
	t, ok := fetchTaskStorage[iString]
	_ = t
	if !ok {
		err = newErrorWithCode(errorcodes.NoFetchTaskToErase, "FetchTask «%s» not found", iString)
		return
	}
	delete(fetchTaskStorage, iString)
	return
}

// It's really a shame to copy the entire contents of map to array,
// but, if we consider a possible use in a concurrent environment,
// it may turn out to be not so bad. And we only copy pointers
func allFetchTasks() []*FetchTask {
	// pre-initialize the result
	result := make([]*FetchTask, len(fetchTaskStorage))
	i := 0
	for _, v := range fetchTaskStorage {
		result[i] = v
		i++
	}
	return result
}
