package errorcodes

//go:generate stringer -type=TaskErrorCode

// TaskErrorCode is returned by service
type TaskErrorCode int

const (
	// NoError means no error, obviously
	NoError TaskErrorCode = iota
	// FailedToParsetaskJSON means an incorrect json in request
	FailedToParsetaskJSON
	// NoTaskToErase means an attempt to delete a non-existent task
	NoTaskToErase
	// FailedToMakeARequest means we were unable to create a request object
	FailedToMakeARequest
	// FailedToSendARequest means we were unable to send a request
	FailedToSendARequest
	// FailedToReadRequestBody only has docs due to linter's wimpyness
	FailedToReadRequestBody
)
