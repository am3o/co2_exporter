package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace      = "tfa_airco2ntrol"
	labelUnit      = "unit"
	labelOperation = "operation"
)

type Collector struct {
	error         *prometheus.GaugeVec
	carbonDioxide *prometheus.GaugeVec
	temperature   *prometheus.GaugeVec
	humidity      *prometheus.GaugeVec
}

func New() *Collector {
	return &Collector{
		error: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:   namespace,
			Name:        "error_total",
			Help:        "Total errors of the measurement",
			ConstLabels: prometheus.Labels{},
		},
			[]string{labelOperation}),
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
	c.error.Collect(metrics)
	c.humidity.Collect(metrics)
	c.temperature.Collect(metrics)
	c.carbonDioxide.Collect(metrics)
}

func (c *Collector) IncError(operation string) {
	c.error.With(prometheus.Labels{labelOperation: operation}).Inc()
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
