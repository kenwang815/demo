package env_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github/demo/env"
)

func Test_GetKeys(t *testing.T) {
	host := "http://localhost"
	port := "8080"

	os.Setenv(env.DBHost, host)
	os.Setenv(env.DBPort, port)

	_v := env.Init()
	keys := _v.Keys()
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, host, _v[env.DBHost])
	assert.Equal(t, port, _v[env.DBPort])
}
