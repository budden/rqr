//go:generate stringer -type=FetchTaskErrorCode

// Package errorcodes contains error codes we return to the client in a json form
package errorcodes

// FetchTaskErrorCode is returned by service
type FetchTaskErrorCode int

const (
	// NoError means no error, obviously
	NoError FetchTaskErrorCode = iota
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
)
