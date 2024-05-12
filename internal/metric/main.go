package metrics

import (
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  //"github.com/prometheus/client_golang/prometheus/promauto"
  "github.com/rs/zerolog"
)

var (
    registry = prometheus.NewRegistry()
    gauges = ()
)


func RegisterAllMetrics(metricsList []string){
    for _, metricName := range metricsList {
	    // Define a new gauge metric
        gauge := prometheus.NewGauge(prometheus.GaugeOpts{
            Name: metricName,
        })
        prometheus.MustRegister(gauge)
	    gauges.AddGaugeToMap(metricName, gauge)
	}
    fmt.Println("gauges" ,gauges)
}

func SetGauge(gaugeName string, gaugeValue float64){

}
