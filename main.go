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

var (
	Version string
)

func main() {
	DevicePath, exists := os.LookupEnv("CO2MONITOR_DEVICE")
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
		logger.With(zap.Error(err)).Fatal("could not open device stream")
	}

	ctx := context.Background()
	defer ctx.Done()

	defer func(ctx context.Context) {
		if err := airController.CloseWithContext(ctx); err != nil {
			logger.With(zap.Error(err)).Fatal("could not close the device connection")
			return
		}

		logger.Info("successfully closed connection to device")
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(10 * time.Second)
		for ; ; <-ticker.C {
			carbonDioxide, temperature, humidity, err := airController.Read(ctx)
			if err != nil {
				logger.With(zap.Error(err)).Error("faulty measurement")
				continue
			}

			collector.SetCarbonDioxideInPPM(carbonDioxide)
			collector.SetTemperatureInCelsius(temperature)
			collector.SetHumidityInPercent(humidity)

			logger.With(
				zap.Float64("carbon_dioxide", carbonDioxide),
				zap.Float64("temperature", temperature),
				zap.Float64("humidity", humidity),
			).Debug("successfully measurement")
		}
	}(ctx)

	http.Handle("/metrics", promhttp.Handler())

	logger.With(
		zap.String("version", Version),
		zap.Int("port", 8080),
	).Info("start co2-exporter service")
	if err := http.ListenAndServe(net.JoinHostPort("", "8080"), nil); err != nil {
		logger.With(zap.Error(err)).Fatal("could not start http service")
	}
}
