package collector

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace      = "tfa_airco2ntrol"
	labelUnit      = "unit"
	labelOperation = "operation"
)

type Collector struct {
	read          prometheus.Histogram
	failure       *prometheus.GaugeVec
	carbonDioxide *prometheus.GaugeVec
	temperature   *prometheus.GaugeVec
	humidity      *prometheus.GaugeVec
}

func New() *Collector {
	return &Collector{
		read: prometheus.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "read_duration_milliseconds",
			Help:      "Histogram of read operation duration",
			Buckets:   []float64{1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144},
		}),
		failure: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "failure_total",
			Help:      "Total number of failed operations",
		}, []string{labelOperation}),
		carbonDioxide: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "carbon_dioxide_total",
			Help:      "Total detected carbon dioxide in ppm of the sensor",
		}, []string{labelUnit}),
		temperature: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "temperature_total",
			Help:      "Total detected temperature of the sensor",
		}, []string{labelUnit}),
		humidity: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "humidity_total",
			Help:      "Total detected humidity of the sensor",
		}, []string{labelUnit}),
	}
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	c.carbonDioxide.Collect(metrics)
	c.failure.Collect(metrics)
	c.humidity.Collect(metrics)
	c.read.Collect(metrics)
	c.temperature.Collect(metrics)
}

func (c *Collector) IncFailure(operation string) {
	c.failure.With(prometheus.Labels{labelOperation: operation}).Inc()
}

func (c *Collector) Track(duration time.Duration) {
	c.read.Observe(float64(duration.Milliseconds()))
}

func (c *Collector) SetCarbonDioxideInPPM(value float64) {
	if value == 0 {
		return
	}

	c.carbonDioxide.With(prometheus.Labels{labelUnit: "ppm"}).Set(value)
}

func (c *Collector) SetTemperatureInCelsius(value float64) {
	if value == 0 {
		return
	}

	c.temperature.With(prometheus.Labels{labelUnit: "celsius"}).Set(value)
}

func (c *Collector) SetHumidityInPercent(value float64) {
	if value == 0 {
		return
	}
	c.humidity.With(prometheus.Labels{labelUnit: "percent"}).Set(value)
}
