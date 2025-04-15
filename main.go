package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/am3o/co2_exporter/pkg/collector"
	"github.com/am3o/co2_exporter/pkg/device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
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

	airController := device.New()
	if err := airController.Open(DevicePath); err != nil {
		logger.Fatal("could not open device stream", zap.Error(err))
	}

	ctx := context.Background()
	defer ctx.Done()

	defer func(ctx context.Context) {
		if err := airController.CloseWithContext(ctx); err != nil {
			logger.Fatal("could not close the device connection", zap.Error(err))
			return
		}

		logger.Info("successfully closed connection to device")
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(10 * time.Second)
		for ; ; <-ticker.C {
			carbonDioxide, temperature, humidity, err := airController.Read(ctx)
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
	}(ctx)

	http.Handle("/internal/metrics", promhttp.Handler())

	logger.Info("start exporter", zap.Int("port", 8080))
	if err := http.ListenAndServe(net.JoinHostPort("", "8080"), nil); err != nil {
		logger.Fatal("could not start http service", zap.Error(err))
	}
}
