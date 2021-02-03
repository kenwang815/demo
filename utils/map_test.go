package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github/demo/utils"
)

type X struct {
	FileUUID string `mapKey:"uuid"`
	FileName string `mapKey:"name,omitempty"`
	Name     string `mapKey:"ignore"`
}

func Test_structToMap_Success(t *testing.T) {
	x := &X{
		FileUUID: "1234",
		FileName: "test",
		Name:     "TestName",
	}
	b := make(map[string]interface{})
	b["uuid"] = "1234"
	b["name"] = "test"
	testMap := utils.Map(x)
	assert.Equal(t, b, testMap)
}

func Test_structToMap_Omitempty(t *testing.T) {
	x := &X{
		FileUUID: "1234",
	}
	b := make(map[string]interface{})
	b["uuid"] = "1234"
	testMap := utils.Map(x)
	assert.Equal(t, b, testMap)
}
