package device

import (
	"context"
	"fmt"
	"io"
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

type Device interface {
	io.ReadCloser
}

type AirController struct {
	device Device
}

func New(path string) (AirController, error) {
	device, err := os.OpenFile(path, os.O_APPEND|os.O_RDONLY, 0)
	if err != nil {
		return AirController{}, fmt.Errorf("could not open file: %w", err)
	}

	enableReportCode := [9]byte{0x0, 0xc4, 0xc6, 0xc0, 0x92, 0x40, 0x23, 0xdc, 0x96}
	_, _, ep := syscall.Syscall(syscall.SYS_IOCTL, device.Fd(), uintptr(enableHidiocsFeature9), uintptr(unsafe.Pointer(&enableReportCode)))
	if ep != 0 {
		return AirController{}, fmt.Errorf("could not enable device to stream values: %w", device.Close())
	}

	return AirController{
		device: device,
	}, nil
}

func (ac *AirController) Read(ctx context.Context) (carbonDioxide float64, temperature float64, humidity float64, err error) {
	for range 4 {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			return
		default:
		}

		buffer := make([]byte, 8)
		_, err = ac.device.Read(buffer)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid measurement detected: %w", err)
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

func (ac *AirController) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return ac.device.Close()
}
