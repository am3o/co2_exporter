package main

import (
	"github.com/Am3o/co2_exporter/pkg/collector"
	"github.com/Am3o/co2_exporter/pkg/device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"os"
	"time"
)

func main() {
	DevicePath, exists := os.LookupEnv("CO2MOINITOR_DEVICE")
	if !exists {
		DevicePath = "/dev/hidraw0"
	}

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	collector := collector.New()
	prometheus.MustRegister(collector)

	var airController device.AirController
	if err := airController.Open(DevicePath); err != nil {
		logger.Error("could not open device stream", zap.Error(err))
		panic(err)
	}
	defer airController.Close()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for ; ; <-ticker.C {
			carbonDioxide, temperature, humidity, err := airController.Read()
			if err != nil {
				logger.Error("faulty measurement", zap.Error(err))
				continue
			}

			collector.SetCarbonDioxideInPPM(carbonDioxide)
			collector.SetTemperatureInCelsius(temperature)
			collector.SetHumidityInPercent(humidity)
			logger.Info("successfully measurement",
				zap.Float64("carbon_dioxide", carbonDioxide),
				zap.Float64("temperature", temperature),
				zap.Float64("humidity", humidity))
		}
	}()

	http.Handle("/internal/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Error("could not start http service", zap.Error(err))
		panic(err)
	}
}
