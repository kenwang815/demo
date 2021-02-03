package service_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github/demo/config"
	"github/demo/daos"
	"github/demo/database"
	"github/demo/model/device"
	"github/demo/service"
	"github/demo/test"
)

type DeviceTestCaseSuite struct {
	db     database.IDatabase
	device service.IDeviceService
}

func setupDeviceTestCaseSuite(t *testing.T) (DeviceTestCaseSuite, func(t *testing.T)) {
	s := DeviceTestCaseSuite{}

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
	s.device = service.NewDeviceService(deviceRepo)

	return s, func(t *testing.T) {
		s.db.Close()
		os.Remove(df.Name())
	}
}

func GetDeviceFromService1() *service.Device {
	return &service.Device{
		Id:      "c9d7c314-fd95-448a-8db9-4756cc774f7d",
		Model:   "Pro",
		Color:   "White",
		Version: "v1.2",
	}
}

func GetDeviceFromService2() *service.Device {
	return &service.Device{
		Id:      "99f970f5-b876-4c94-9190-34ee11d54edb",
		Model:   "Normal",
		Color:   "Black",
		Version: "v1.2",
	}
}

func GetDeviceFromService3() *service.Device {
	return &service.Device{
		Id:      "c9d7c314-fd95-448a-8db9-4756cc774f7d",
		Version: "v1.6",
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

func TestDeviceService_Find(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description   string
		page          *service.Page
		filter        *service.Device
		expectedCode  service.ErrorCode
		expected      []*service.Device
		setupTestCase test.SetupSubTest
	}{
		{
			description:  "input page=1 number=2",
			page:         &service.Page{Page: 1, Number: 2},
			filter:       &service.Device{},
			expectedCode: service.ErrorCodeSuccess,
			expected: []*service.Device{
				GetDeviceFromService1(),
				GetDeviceFromService2(),
			},
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())
				s.db.GetDB().Create(GetDevice2())

				return func(t *testing.T) {
					s.db.GetDB().DropTable(&device.Device{})
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			devices, code := s.device.Find(tc.filter, tc.page)
			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expected, devices)
		})
	}
}

func TestDeviceService_Register(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description   string
		inputData     *service.Device
		expectedCode  service.ErrorCode
		expected      *service.Device
		setupTestCase test.SetupSubTest
	}{
		{
			description:  "success",
			inputData:    GetDeviceFromService1(),
			expectedCode: service.ErrorCodeSuccess,
			expected:     GetDeviceFromService1(),
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().AutoMigrate(&device.Device{})

				return func(t *testing.T) {
					s.db.GetDB().DropTable(&device.Device{})
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			device, code := s.device.Register(tc.inputData)
			assert.Equal(t, tc.expectedCode, code)

			if tc.expected != nil {
				tc.expected.Id = device.Id
				tc.expected.CreateTime = device.CreateTime
				tc.expected.UpdateTime = device.UpdateTime
			}
			assert.Equal(t, tc.expected, device)

		})
	}
}

func TestDeviceService_Update(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description    string
		inputData      *service.Device
		expectedCode   service.ErrorCode
		expectedAffect int64
		setupTestCase  test.SetupSubTest
	}{
		{
			description:    "success",
			inputData:      GetDeviceFromService3(),
			expectedCode:   service.ErrorCodeSuccess,
			expectedAffect: 1,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
					s.db.GetDB().DropTable(&device.Device{})
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)
			affect, code := s.device.Update(tc.inputData)

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedAffect, affect)
		})
	}
}

func TestDeviceService_Delete(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		description   string
		inputId       string
		expectedCode  service.ErrorCode
		setupTestCase test.SetupSubTest
	}{
		{
			description:  "success",
			inputId:      GetDeviceFromService1().Id,
			expectedCode: service.ErrorCodeSuccess,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
					s.db.GetDB().DropTable(&device.Device{})
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.description, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			code := s.device.Delete(tc.inputId)
			assert.Equal(t, tc.expectedCode, code)
		})
	}
}
