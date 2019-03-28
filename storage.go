package main

import (
	"math/big"

	"github.com/budden/rqr/pkg/errorcodes"
)

var queryID = big.NewInt(0)
var one = big.NewInt(1)

var taskStorage map[string]*Task

func saveTask(pt *ParsedTask, et *ExecutedTask) (t *Task) {
	// FIXME will the map be of size 1?
	queryID.Add(queryID, one)
	iString := queryID.String()
	taskStorage[iString] = &Task{ID: iString, pt: pt, et: et}
	return
}

func getTask(iString string) (t *Task, ok bool) {
	t, ok = taskStorage[iString]
	return
}

func eraseTask(iString string) (err *errorWithCode) {
	t, ok := taskStorage[iString]
	_ = t
	if !ok {
		err = newErrorWithCode(errorcodes.NoTaskToErase, "Task «%s» not found", iString)
		return
	}
	delete(taskStorage, iString)
	return
}

// It's really a shame to copy the entire contents of map to array,
// but, if we consider a possible use in a concurrent environment,
// it may turn out to be not so bad. And we only copy pointers
func allTasks() []*Task {
	// pre-initialize the result
	result := make([]*Task, len(taskStorage))
	i := 0
	for _, v := range taskStorage {
		result[i] = v
		i++
	}
	return result
}
