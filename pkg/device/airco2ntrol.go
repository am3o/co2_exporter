package device

import (
	"errors"
	"fmt"
	"github.com/am3o/co2_exporter/pkg/model"
	"os"
)

const (
	OperationCarbonDioxide = 'P'
	OperationTemperature   = 'B'
	OperationHumidity      = 'A'
)

type AirController struct {
	device *os.File
}

func (ac *AirController) Open(path string) (err error) {
	ac.device, err = os.OpenFile(path, os.O_APPEND|os.O_RDONLY, 0)
	if err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

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
