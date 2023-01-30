package device

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/am3o/co2_exporter/pkg/model"
)

const (
	OperationCarbonDioxide = 'P'
	OperationTemperature   = 'B'
	OperationHumidity      = 'A'

	enableHidiocsFeature9 = 0xC0094806
)

type AirController struct {
	device           *os.File
	enableReportCode unsafe.ArbitraryType
}

func New() AirController {
	return AirController{
		enableReportCode: *unsafe.Pointer(&[9]byte{0x0, 0xc4, 0xc6, 0xc0, 0x92, 0x40, 0x23, 0xdc, 0x96}),
	}
}

func (ac *AirController) Open(path string) error {
	device, err := os.OpenFile(path, os.O_APPEND|os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, device.Fd(), uintptr(enableHidiocsFeature9), uintptr(ac.enableReportCode))
	if ep != 0 {
		return fmt.Errorf("could not enable device to stream values: %w", device.Close())
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
