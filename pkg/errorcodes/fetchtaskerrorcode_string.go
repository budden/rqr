// Code generated by "stringer -type=FetchTaskErrorCode"; DO NOT EDIT.

package errorcodes

import "strconv"

const _FetchTaskErrorCode_name = "OKFailedToParsefetchTaskJSONFetchTaskNotFoundNoFetchTaskToEraseFailedToMakeARequestFailedToSendARequestFailedToReadRequestBodyDuplicateParamInTheFormIncorrectIDFormatIncorrectRequestMethodIncorrectURLUnknownError"

var _FetchTaskErrorCode_index = [...]uint8{0, 2, 28, 45, 63, 83, 103, 126, 149, 166, 188, 200, 212}

func (i FetchTaskErrorCode) String() string {
	if i < 0 || i >= FetchTaskErrorCode(len(_FetchTaskErrorCode_index)-1) {
		return "FetchTaskErrorCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _FetchTaskErrorCode_name[_FetchTaskErrorCode_index[i]:_FetchTaskErrorCode_index[i+1]]
}
