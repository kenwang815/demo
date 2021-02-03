/**
xeniro platform custom error code
http status code + category + serial number
xxx 00 xx General
xxx 01 xx MdcSrv
*/

package service

import (
	"fmt"
	"reflect"
	"strconv"
)

type ErrorCode int

func (e ErrorCode) Int() int     { return int(e) }
func (e ErrorCode) Int64() int64 { return int64(e) }

// 200 00
const (
	ErrorCodeSuccess ErrorCode = iota + 2000000
	ErrorCodeSuccessButNotFound
)

// 400 00
const (
	ErrorCodeBadRequest ErrorCode = iota + 4000000
	ErrorCodeParseUUIDFail
)

// 401 00
const (
	ErrorCodeTokenInvalid ErrorCode = iota + 4010000
	ErrorCodeTokenExpired
)

// 403 00
const (
	ErrorCodeForbidden ErrorCode = iota + 4030000
)

// 404 00
const (
	ErrorCodeNotFound ErrorCode = iota + 4040000
)

// 500 00
const (
	ErrorCodeServerErr ErrorCode = iota + 5000000
	ErrorCodeDatabaseFail
	ErrorCodeTokenCreateFail
)

// 500 01
const (
	ErrorCodeDeviceDBFindFail ErrorCode = iota + 5000100
	ErrorCodeDeviceDBUpdateFail
	ErrorCodeDeviceDBCreateFail
	ErrorCodeDeviceDBDeleteFail
)

var errorMsg = map[ErrorCode]string{
	ErrorCodeSuccess:            "Success",
	ErrorCodeSuccessButNotFound: "Success with no affect rows",
	ErrorCodeBadRequest:         "Bad request",
	ErrorCodeTokenInvalid:       "Token invalid",
	ErrorCodeForbidden:          "Forbidden",
	ErrorCodeNotFound:           "Not found",
	ErrorCodeServerErr:          "Internal server error",
	ErrorCodeDatabaseFail:       "Database failure",
	ErrorCodeTokenCreateFail:    "Token create fail",
	ErrorCodeParseUUIDFail:      "Had a error in uuid parsing",
	ErrorCodeDeviceDBFindFail:   "Device find fail",
	ErrorCodeDeviceDBUpdateFail: "Device update fail",
	ErrorCodeDeviceDBCreateFail: "Device create fail",
	ErrorCodeDeviceDBDeleteFail: "Device delete fail",
}

func ErrorMsg(code ErrorCode) string {
	return errorMsg[code]
}

func ErrorStatusCode(code ErrorCode) int {
	i, _ := strconv.Atoi(fmt.Sprintf("%d", code)[:3])
	return i
}

func NewErrors(s ...interface{}) *errors {
	r := &errors{
		[]interface{}{},
	}

	for _, v := range s {
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Slice, reflect.Array:
			v1 := reflect.ValueOf(v)
			for i := 0; i < v1.Len(); i++ {
				r.s = append(r.s, v1.Index(i).Interface())
			}
		default:
			r.s = append(r.s, v)
		}
	}

	return r
}

type errors struct {
	s []interface{}
}

func (m *errors) Add(s interface{}) *errors {
	m.s = append(m.s, s)
	return m
}

func (m *errors) Error() []interface{} {
	return m.s
}

func (m *errors) NotEmpty() bool {
	return len(m.s) > 0
}
