package errorcodes

//go:generate stringer -type=InquiryErrorCode

// InquiryErrorCode is returned by service
type InquiryErrorCode int

const (
	// NoError means no error, obviously
	NoError InquiryErrorCode = iota
	// FailedToParseInquiryJson means incorrect json in request
	FailedToParseInquiryJson
)
