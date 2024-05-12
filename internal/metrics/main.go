package metrics

import (
  "github.com/prometheus/client_golang/prometheus"
  "github.com/prometheus/client_golang/prometheus/promhttp"
  //"github.com/prometheus/client_golang/prometheus/promauto"
  //"github.com/rs/zerolog"
  "net/http"
)

var (
    registry = prometheus.NewRegistry()
    gauges = NewGaugeMap()
    promHttpHandler = promhttp.Handler()
)

func GetPromHttp() (http.Handler){
    return promHttpHandler
}

func RegisterAllMetrics(metricsList []string){
    for _, metricName := range metricsList {
	    if metricName != "" {
       	    // Define a new gauge metric
            gauge := prometheus.NewGauge(prometheus.GaugeOpts{
                Name: metricName,
            })
            prometheus.MustRegister(gauge)
            gauges.AddGaugeToMap(metricName, gauge)
	    }
	}
}

func SetGauge(gaugeName string, gaugeValue float64){

}
