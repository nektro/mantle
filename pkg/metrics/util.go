package metrics

import (
	"github.com/nektro/mantle/pkg/db"

	"github.com/prometheus/client_golang/prometheus"
)

func newGauge(name string) prometheus.Gauge {
	r := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "astheno_mantle",
			Name:      name,
		},
	)
	prometheus.MustRegister(r)
	return r
}

func newGaugeLabeled(name string, tags ...string) *prometheus.GaugeVec {
	r := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "astheno_mantle",
			Name:      name,
		},
		tags,
	)
	prometheus.MustRegister(r)
	return r
}

func getPropInt(name string) float64 {
	return float64(db.Props.GetInt64(name))
}
