package service

type IDeviceService interface {
	Find(*Device, *Page) ([]*Device, ErrorCode)
	Register(*Device) (*Device, ErrorCode)
	Update(*Device) (int64, ErrorCode)
	Delete(string) ErrorCode
}
