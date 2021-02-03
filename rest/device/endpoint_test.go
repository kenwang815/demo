package device_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github/demo/config"
	"github/demo/daos"
	"github/demo/database"
	"github/demo/model/device"
	routeDevice "github/demo/rest/device"
	"github/demo/service"
	"github/demo/test"
)

type DeviceTestCaseSuite struct {
	db database.IDatabase
	c  *gin.Engine
}

func setupDeviceTestCaseSuite(t *testing.T) (DeviceTestCaseSuite, func(t *testing.T)) {
	s := DeviceTestCaseSuite{
		c: gin.New(),
	}
	s.c.Use(gin.Recovery())

	name := filepath.Join(os.TempDir(), "gorm"+uuid.Must(uuid.NewV4()).String()+".db")
	df, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if df == nil || err != nil {
		panic(fmt.Sprintf("No error should happen when creating db file, but got %+v", err))
	}

	c := &config.Database{
		Dialect: "sqlite",
		Host:    df.Name(),
	}

	s.db, err = database.NewDatabase(c)
	deviceRepo := daos.NewDeviceRepo(s.db.GetDB())
	service.DeviceService = service.NewDeviceService(deviceRepo)

	routeDevice.MakeHandler(s.c.Group("/v1"))

	return s, func(t *testing.T) {
		s.db.Close()
		os.Remove(df.Name())
	}
}

func GetDevice1() *device.Device {
	return &device.Device{
		Id:      "c9d7c314-fd95-448a-8db9-4756cc774f7d",
		Model:   "Pro",
		Color:   "White",
		Version: "v1.2",
	}
}

func GetDevice2() *device.Device {
	return &device.Device{
		Id:      "99f970f5-b876-4c94-9190-34ee11d54edb",
		Model:   "Normal",
		Color:   "Black",
		Version: "v1.2",
	}
}

func TestDeviceFindHandler(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		params       map[string]string
		expected     string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description: "input page=1, number=2",
			route:       "/v1/device",
			method:      "GET",
			params: map[string]string{
				"page":   "1",
				"number": "2",
			},
			expected:     `{"code":2000000,"data":{"datas":[{"Id":"c9d7c314-fd95-448a-8db9-4756cc774f7d","Model":"Pro","Color":"White","Version":"v1.2","CreateTime":0,"UpdateTime":0},{"Id":"99f970f5-b876-4c94-9190-34ee11d54edb","Model":"Normal","Color":"Black","Version":"v1.2","CreateTime":0,"UpdateTime":0}],"page":{"page":1,"number":2}},"msg":"Success"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())
				s.db.GetDB().Create(GetDevice2())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, nil)
			q := req.URL.Query()
			for k, v := range tc.params {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)
			assert.Equal(t, tc.expectedCode, actul.Code)
			assert.Equal(t, tc.expected, strings.Replace(actul.Body.String(), "\n", "", -1))
		})
	}
}

func TestDeviceRegisterHandler(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		body         string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description:  "success",
			route:        "/v1/device",
			method:       "POST",
			body:         `{"model":"Pro","color":"green","version":"2.0"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)

			assert.Equal(t, tc.expectedCode, actul.Code)
		})
	}
}

func TestDeviceDeleteHandler(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description:  "success",
			route:        "/v1/device/" + GetDevice1().Id.String(),
			method:       "DELETE",
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, nil)
			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)

			assert.Equal(t, tc.expectedCode, actul.Code)
		})
	}
}

func TestDeviceUpdateHandler(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description  string
		route        string
		method       string
		body         string
		expectedCode int
		setupSubTest test.SetupSubTest
	}{
		{
			description:  "success",
			route:        "/v1/device",
			method:       "PUT",
			body:         `{"id":"c9d7c314-fd95-448a-8db9-4756cc774f7d","color":"White"}`,
			expectedCode: http.StatusOK,
			setupSubTest: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupSubTest(t)
			defer teardownSubTest(t)

			req := httptest.NewRequest(tc.method, tc.route, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", gin.MIMEJSON)
			actul := httptest.NewRecorder()
			s.c.ServeHTTP(actul, req)
			assert.Equal(t, tc.expectedCode, actul.Code)
		})
	}
}
