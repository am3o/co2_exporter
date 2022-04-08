package main

import (
	"github.com/Am3o/co2_exporter/pkg/collector"
	"github.com/Am3o/co2_exporter/pkg/device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

const (
	DevicePath = "/dev/hidraw0"
)

func main() {
	collector := collector.New()
	prometheus.MustRegister(collector)

	var airController device.AirController
	if err := airController.Open(DevicePath); err != nil {
		panic(err)
	}
	defer airController.Close()

	go func() {
		for range time.NewTicker(5 * time.Second).C {
			carbonDioxide, temperature, humidity, err := airController.Read()
			if err != nil {
				continue
			}

			collector.SetCarbonDioxideInPPM(carbonDioxide)
			collector.SetTemperatureInCelsius(temperature)
			collector.SetHumidityInPercent(humidity)
		}
	}()

	http.Handle("/internal/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
