package device

import (
	"errors"
	"fmt"
	"github.com/Am3o/co2_exporter/pkg/model"
	"os"
	"syscall"
	"unsafe"
)

const (
	OperationCarbonDioxide = 'P'
	OperationTemperature   = 'B'
	OperationHumidity      = 'A'
)

type AirController struct {
	flag   [9]byte
	device *os.File
}

func NewAirController() AirController {
	return AirController{
		flag: [9]byte{0x0, 0xc4, 0xc6, 0xc0, 0x92, 0x40, 0x23, 0xdc, 0x96},
	}
}

func (ac *AirController) Open(path string) error {
	device, err := os.OpenFile(path, os.O_APPEND|os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	hidiocsFeature := int64(0xC0094806)
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, device.Fd(), uintptr(hidiocsFeature), uintptr(unsafe.Pointer(&ac.flag)))
	if ep != 0 {
		return device.Close()
	}

	ac.device = device
	return nil
}

func (ac *AirController) Read() (carbonDioxide float64, temperature float64, humidity float64, err error) {
	for i := 0; i < 4; i++ {
		buffer := make([]byte, 8)
		_, err = ac.device.Read(buffer)
		if err != nil {
			return 0, 0, 0, errors.New("invalid measurement detected")
		}

		signal := model.Signal(buffer)
		switch signal.Type() {
		case OperationCarbonDioxide:
			carbonDioxide = signal.Value()
		case OperationTemperature:
			temperature = signal.Value()/16 - 273.15
		case OperationHumidity:
			humidity = signal.Value() / 100
		default:
			continue
		}
	}
	return
}

func (ac *AirController) Close() error {
	return ac.device.Close()
}
