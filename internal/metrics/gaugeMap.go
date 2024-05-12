package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
)


type GaugeMap struct {
    gauges map[string]prometheus.Gauge
}

// NewGaugeMap initializes a new GaugeMap.
func NewGaugeMap() *GaugeMap {
    return &GaugeMap{
        gauges: make(map[string]prometheus.Gauge),
    }
}

// RegisterGauge registers a new gauge with the given name and help string.
func (gm *GaugeMap) AddGaugeToMap(name string , gauge prometheus.Gauge) {
    gm.gauges[name] = gauge
}

// Set sets the value of the gauge with the given name.
func (gm *GaugeMap) Set(name string, value float64) {
    if gauge, ok := gm.gauges[name]; ok {
        gauge.Set(value)
    }
}

// Get retrieves the value of the gauge with the given name.
//func (gm *GaugeMap) Get(name string) float64 {
//    if gauge, ok := gm.gauges[name]; ok {
//        return gauge
//    }
//    return -1
//}