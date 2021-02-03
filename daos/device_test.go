package daos_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"

	"github/demo/config"
	"github/demo/daos"
	"github/demo/database"
	"github/demo/model"
	"github/demo/model/device"
	"github/demo/test"
)

type DeviceTestCaseSuite struct {
	db         database.IDatabase
	deviceRepo device.Repository
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
	s.deviceRepo = daos.NewDeviceRepo(s.db.GetDB())
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

func UpdateDevice1() *device.Device {
	return &device.Device{
		Id:      "c9d7c314-fd95-448a-8db9-4756cc774f7d",
		Model:   "Pro",
		Color:   "White",
		Version: "v1.5",
	}
}

func UpdateDevice2() *device.Device {
	return &device.Device{
		Id: "c9d7c314-fd95-448a-8db9-4756cc774f7d",
	}
}

func TestDeviceDaos_Get(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		id            device.UUID
		wantResult    *device.Device
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "success",
			id:         GetDevice1().Id,
			wantResult: GetDevice1(),
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "no data",
			id:         GetDevice1().Id,
			wantResult: nil,
			err:        errors.New("record not found"),
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			p, err := s.deviceRepo.Get(tc.id)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, p, tc.wantResult)
			}
		})
	}
}

func TestDeviceDaos_Create(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *device.Device
		wantResult    *device.Device
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "success",
			testData:   GetDevice1(),
			wantResult: GetDevice1(),
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			p, err := s.deviceRepo.Create(tc.testData)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				tc.wantResult.CreateTime = p.CreateTime
				assert.Equal(t, p, tc.wantResult)
			}
		})
	}
}

func TestDeviceDaos_Update(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *device.Device
		wantResult    *device.Device
		rowAffected   int64
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:        "success",
			testData:    UpdateDevice1(),
			wantResult:  UpdateDevice1(),
			rowAffected: 1,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:        "ignore_uuid_update",
			testData:    UpdateDevice2(),
			wantResult:  nil,
			rowAffected: 0,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			device, affected, err := s.deviceRepo.Update(tc.testData)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.rowAffected, affected)
			if err == nil && device != nil {
				tc.wantResult.UpdateTime = device.UpdateTime
				d, _ := s.deviceRepo.Get(tc.testData.Id)
				assert.Equal(t, tc.testData, d)
			}
		})
	}
}

func TestDeviceDaos_Delete(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		id            device.UUID
		rowAffected   int64
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:        "success",
			id:          GetDevice1().Id,
			rowAffected: 1,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:        "not exist id",
			id:          GetDevice2().Id,
			rowAffected: 0,
			err:         nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())

				return func(t *testing.T) {
				}
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			affected, err := s.deviceRepo.Delete(tc.id)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.rowAffected, affected)
		})
	}
}

func TestDeviceDaos_List(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *device.Device
		wantResult    []*device.Device
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "no data",
			testData:   &device.Device{},
			wantResult: []*device.Device{},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "get all data",
			testData:   &device.Device{},
			wantResult: []*device.Device{GetDevice1(), GetDevice2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())
				s.db.GetDB().Create(GetDevice2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "find id",
			testData:   &device.Device{Id: GetDevice1().Id},
			wantResult: []*device.Device{GetDevice1()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			devices, err := s.deviceRepo.List(tc.testData)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, devices, tc.wantResult)
			}
		})
	}
}

func TestDeviceDaos_Find(t *testing.T) {
	s, teardownTestCase := setupDeviceTestCaseSuite(t)
	defer teardownTestCase(t)

	tt := []struct {
		name          string
		testData      *device.Device
		testPage      *model.Page
		wantResult    []*device.Device
		err           error
		setupTestCase test.SetupSubTest
	}{
		{
			name:       "no data",
			testData:   &device.Device{},
			testPage:   &model.Page{Limit: 0, Offset: 0},
			wantResult: []*device.Device{},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())
				s.db.GetDB().Create(GetDevice2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "input limit > count",
			testData:   &device.Device{},
			testPage:   &model.Page{Limit: 3, Offset: 0},
			wantResult: []*device.Device{GetDevice1(), GetDevice2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
				s.db.GetDB().DropTable(&device.Device{})
				s.db.GetDB().AutoMigrate(&device.Device{})
				s.db.GetDB().Create(GetDevice1())
				s.db.GetDB().Create(GetDevice2())

				return func(t *testing.T) {
				}
			},
		},
		{
			name:       "find id",
			testData:   &device.Device{Id: GetDevice2().Id},
			testPage:   &model.Page{Limit: 2, Offset: 0},
			wantResult: []*device.Device{GetDevice2()},
			err:        nil,
			setupTestCase: func(t *testing.T) func(t *testing.T) {
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
		t.Run(tc.name, func(t *testing.T) {
			teardownSubTest := tc.setupTestCase(t)
			defer teardownSubTest(t)

			devices, err := s.deviceRepo.Find(tc.testData, tc.testPage)
			if err != nil {
				assert.EqualError(t, err, tc.err.Error(), "An error was expected")
			} else {
				assert.Equal(t, devices, tc.wantResult)
			}
		})
	}
}
