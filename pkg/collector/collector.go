package collector

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "tfa_airco2ntrol"
	labelUnit = "unit"
)

type Collector struct {
	CarbonDioxideGauge *prometheus.GaugeVec
	TemperatureGauge   *prometheus.GaugeVec
	HumidityGauge      *prometheus.GaugeVec
}

func New() *Collector {
	return &Collector{
		CarbonDioxideGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "carbon_dioxide_total",
			Help:      "",
		}, []string{labelUnit}),
		TemperatureGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "temperature_total",
			Help:      "",
		}, []string{labelUnit}),
		HumidityGauge: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "humidity_total",
			Help:      "",
		}, []string{labelUnit}),
	}
}

func (c *Collector) Describe(descs chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, descs)
}

func (c *Collector) Collect(metrics chan<- prometheus.Metric) {
	c.TemperatureGauge.Collect(metrics)
	c.HumidityGauge.Collect(metrics)
	c.CarbonDioxideGauge.Collect(metrics)
}

func (c *Collector) SetCarbonDioxideInPPM(value float64) {
	if value == 0 {
		return
	}

	c.CarbonDioxideGauge.With(prometheus.Labels{labelUnit: "ppm"}).Set(value)
}

func (c *Collector) SetTemperatureInCelsius(value float64) {
	if value == 0 {
		return
	}

	c.TemperatureGauge.With(prometheus.Labels{labelUnit: "celsius"}).Set(value)
}

func (c *Collector) SetHumidityInPercent(value float64) {
	if value == 0 {
		return
	}
	c.HumidityGauge.With(prometheus.Labels{labelUnit: "percent"}).Set(value)
}
