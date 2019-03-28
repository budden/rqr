package main

import (
	"fmt"
	"math/big"

	"github.com/budden/rqr/pkg/errorcodes"
)

var queryID = big.NewInt(0)
var one = big.NewInt(1)

var taskStorage map[string]*Task

func saveTask(pi *ParsedTask, ei *ExecutedTask) (t *Task) {
	// FIXME will the map be of size 1?
	queryID.Add(queryID, one)
	iString := queryID.String()
	taskStorage[iString] = &Task{ID: iString, pi: pi, ei: ei}
	return
}

func getTask(iString string) (t *Task, ok bool) {
	t, ok = taskStorage[iString]
	return
}

func eraseTask(iString string) (err *jsonError) {
	t, ok := taskStorage[iString]
	_ = t
	if !ok {
		msg := fmt.Sprintf("Task «%s» not found", iString)
		err = &jsonError{Code: errorcodes.NoTaskToErase, Message: msg}
		return
	}
	delete(taskStorage, iString)
	return
}
