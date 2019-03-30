//go:generate stringer -type=FetchTaskErrorCode

// Package errorcodes contains error codes we return to the client in a json form
package errorcodes

// FetchTaskErrorCode is returned by service
type FetchTaskErrorCode int

const (
	// OK means no error, obviously
	OK FetchTaskErrorCode = iota
	// FailedToParsefetchTaskJSON means an incorrect json in request
	FailedToParsefetchTaskJSON
	// NoFetchTaskToErase means an attempt to delete a non-existent fetchTask
	NoFetchTaskToErase
	// FailedToMakeARequest means we were unable to create a request object
	FailedToMakeARequest
	// FailedToSendARequest means we were unable to send a request
	FailedToSendARequest
	// FailedToReadRequestBody only has docs due to linter's wimpyness
	FailedToReadRequestBody
	// DuplicateParamInTheForm means that there is more than one parameter with this name
	DuplicateParamInTheForm
	// IncorrectIDFormat is signaled if fetchtaskget is invoked for a non-numeric ID
	IncorrectIDFormat
	// IncorrectRequestMethod means GET instead of POST or vice versa
	IncorrectRequestMethod
	// UnknownError ...
	UnknownError
)
