package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/am3o/co2_exporter/pkg/collector"
	"github.com/am3o/co2_exporter/pkg/device"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	Version string
)

func main() {
	DevicePath, exists := os.LookupEnv("CO2MONITOR_DEVICE")
	if !exists {
		DevicePath = "/dev/hidraw0"
	}

	collector := collector.New()
	prometheus.MustRegister(collector)

	ctx := context.Background()
	defer ctx.Done()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	airController, err := device.New(DevicePath)
	if err != nil {
		logger.
			With(slog.Any("error", err)).
			ErrorContext(ctx, "could not open device stream")
		os.Exit(1)
	}

	defer func(ctx context.Context) {
		if err := airController.Close(ctx); err != nil {
			logger.
				With(slog.Any("error", err)).
				ErrorContext(ctx, "could not close the device connection")
			os.Exit(1)
		}

		logger.InfoContext(ctx, "successfully closed connection to device")
	}(ctx)

	go func(ctx context.Context) {
		ticker := time.NewTicker(10 * time.Second)
		for ; ; <-ticker.C {
			carbonDioxide, temperature, humidity, err := airController.Read(ctx)
			if err != nil {
				logger.
					With(slog.Any("error", err)).
					ErrorContext(ctx, "faulty measurement")
				continue
			}

			collector.SetCarbonDioxideInPPM(carbonDioxide)
			collector.SetTemperatureInCelsius(temperature)
			collector.SetHumidityInPercent(humidity)

			logger.With(
				slog.Float64("carbon_dioxide", carbonDioxide),
				slog.Float64("temperature", temperature),
				slog.Float64("humidity", humidity),
			).DebugContext(ctx, "successfully measurement")
		}
	}(ctx)

	http.Handle("/metrics", promhttp.Handler())

	logger.With(
		slog.String("version", Version),
		slog.Int("port", 8080),
	).Info("start co2-exporter service")
	if err := http.ListenAndServe(net.JoinHostPort("", "8080"), nil); err != nil {
		logger.
			With(slog.Any("error", err)).
			ErrorContext(ctx, "could not start http service")
		os.Exit(1)
	}
}
