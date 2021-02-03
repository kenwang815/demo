package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github/demo/service"
)

// test error status code
func Test_ErrorStatusCode(t *testing.T) {
	statsCode := service.ErrorStatusCode(service.ErrorCodeNotFound)
	assert.Equal(t, statsCode, 404)
}

// test error message
func Test_ErrorMsg(t *testing.T) {
	errMsg := service.ErrorMsg(service.ErrorCodeTokenInvalid)
	assert.Equal(t, errMsg, "Token invalid")
}

// test new detail errors
func Test_NewErrors(t *testing.T) {
	err := service.NewErrors()
	assert.Equal(t, err.Error(), []interface{}{})
	assert.Equal(t, err.NotEmpty(), false)

	err = service.NewErrors("this is error message")
	assert.Equal(t, err.Error(), []interface{}{"this is error message"})
	assert.Equal(t, err.NotEmpty(), true)

	err = service.NewErrors()
	err.Add("error message 1")
	err.Add("error message 2")
	assert.Equal(t, err.Error(), []interface{}{"error message 1", "error message 2"})
	assert.Equal(t, err.NotEmpty(), true)
}
