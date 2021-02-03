package env

import (
	"os"
	"reflect"
)

const (
	LogEnv     = "Logger_Env"
	LogLevel   = "Logger_Level"
	LogFile    = "Logger_Filename"
	DBDialect  = "DB_Dialect"
	DBHost     = "DB_Host"
	DBPort     = "DB_Port"
	DBName     = "DB_Name"
	DBUser     = "DB_User"
	DBPassword = "DB_Password"
)

var eVar []string = []string{
	LogEnv,
	LogLevel,
	LogFile,
	DBDialect,
	DBHost,
	DBPort,
	DBName,
	DBUser,
	DBPassword,
}

type Variables map[string]interface{}

func (v *Variables) Keys() (keys []string) {
	value := reflect.ValueOf(v).Elem()
	ks := value.MapKeys()
	for _, key := range ks {
		keys = append(keys, key.String())
	}
	return
}

func Init() Variables {
	vars := make(Variables)
	for _, p := range eVar {
		v := os.Getenv(p)
		if len(v) > 0 {
			vars[p] = v
		}
	}

	return vars
}
