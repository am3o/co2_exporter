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

	ctx := context.Background()
	defer ctx.Done()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.
		With(slog.String("version", Version)).
		InfoContext(ctx, "airco2ntrol CO₂ Exporter")

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

	collector := collector.New()
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())

	go func(ctx context.Context) {
		ticker := time.NewTicker(15 * time.Second)
		for ; ; <-ticker.C {
			now := time.Now()
			carbonDioxide, temperature, humidity, err := airController.Read(ctx)
			if err != nil {
				logger.
					With(slog.Any("error", err)).
					ErrorContext(ctx, "faulty measurement")
				collector.IncFailure("read_sensor")
				collector.Track(time.Since(now))
				continue
			}

			collector.SetCarbonDioxideInPPM(carbonDioxide)
			collector.SetTemperatureInCelsius(temperature)
			collector.SetHumidityInPercent(humidity)
			collector.Track(time.Since(now))
		}
	}(ctx)

	logger.
		With(slog.Int("port", 8080)).
		InfoContext(ctx, "start co2-exporter service")

	server := new(http.Server{
		Addr:              net.JoinHostPort("", "8080"),
		ReadHeaderTimeout: 3 * time.Second,
	})
	if err := server.ListenAndServe(); err != nil {
		logger.
			With(slog.Any("error", err)).
			ErrorContext(ctx, "could not start http service")
		os.Exit(1)
	}
}
