package main

import (
	"math/big"
	"sync"

	"github.com/budden/rqr/pkg/errorcodes"
)

var queryID = big.NewInt(0)
var one = big.NewInt(1)

var fetchTaskStorage = map[string]*FetchTask{}

// must ensure that all operations are protected by
// a mutex because request handlers run in multiple goroutines
var fetchTaskStorageMutex sync.Mutex

func saveFetchTask(pt *ParsedFetchTask, et *ExecutedFetchTask) (t *FetchTask) {
	fetchTaskStorageMutex.Lock()
	defer fetchTaskStorageMutex.Unlock()
	queryID.Add(queryID, one)
	thisID := &big.Int{}
	thisID.Set(queryID)
	iString := queryID.String()
	t = &FetchTask{ID: iString, IDn: thisID, pt: pt, et: et}
	fetchTaskStorage[iString] = t
	return
}

func getFetchTask(iString string) (t *FetchTask, ok bool) {
	fetchTaskStorageMutex.Lock()
	defer fetchTaskStorageMutex.Unlock()
	t, ok = fetchTaskStorage[iString]
	return
}

func eraseFetchTask(iString string) (err ErrorWithCode) {
	fetchTaskStorageMutex.Lock()
	defer fetchTaskStorageMutex.Unlock()
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
	fetchTaskStorageMutex.Lock()
	defer fetchTaskStorageMutex.Unlock()
	// pre-initialize the result
	result := make([]*FetchTask, len(fetchTaskStorage))
	i := 0
	for _, v := range fetchTaskStorage {
		result[i] = v
		i++
	}
	return result
}
