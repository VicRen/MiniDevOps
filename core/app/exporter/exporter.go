package exporter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/vicren/covid-away/core/common"
	"github.com/vicren/covid-away/core/features/exporter"
)

type Handler struct {
	config        *Config
	cLock         sync.RWMutex
	counterMap    map[string]prometheus.Counter
	gaugeMap      map[string]prometheus.Gauge
	histogramMap  map[string]prometheus.Histogram
	summaryMap    map[string]prometheus.Summary
	summaryVecMap map[string]*prometheus.SummaryVec
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	h := &Handler{
		counterMap: make(map[string]prometheus.Counter),
		gaugeMap:   make(map[string]prometheus.Gauge),
	}
	h.config = config
	return h, nil
}

func (h *Handler) Type() interface{} {
	return exporter.HandlerType()
}

func (h *Handler) Start() error {
	fmt.Println("----->Exporter start")
	http.Handle("/metrics", promhttp.Handler())
	if len(h.config.Address) > 0 {
		go func() {
			fmt.Println("----->listing on:", h.config.Address)
			if err := http.ListenAndServe(h.config.Address, nil); err != nil {
				panic(err)
			}
		}()
		return nil
	}
	return fmt.Errorf("invalid address: %s", h.config.Address)
}

func (h *Handler) Close() error {
	fmt.Println("----->Exporter close")
	return nil
}

func (h *Handler) CounterInc(name string, labels map[string]string) error {
	newName := name + labelString(labels)
	c, ok := h.counterMap[newName]
	if ok {
		c.Inc()
		return nil
	}
	c = promauto.NewCounter(prometheus.CounterOpts{
		Name:        name,
		Help:        "testing-help",
		ConstLabels: labels,
	})
	h.counterMap[newName] = c
	c.Inc()
	return nil
}

func (h *Handler) CounterAdd(name string, labels map[string]string) error {
	return nil
}

func (h *Handler) GaugeSet(name string, labels map[string]string, value float64) error {
	newName := name + labelString(labels)
	c, ok := h.gaugeMap[newName]
	if ok {
		c.Set(value)
		return nil
	}
	c = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        name,
		Help:        "testing-gauge-help",
		ConstLabels: labels,
	})
	h.gaugeMap[newName] = c
	c.Set(value)
	return nil
}

func (h *Handler) HistogramObserve(name string, labels map[string]string, value float64) error {
	return nil
}

func (h *Handler) SummaryObserve(name string, labels map[string]string, value float64) error {
	return nil
}

func labelString(labels map[string]string) string {
	buf, err := json.Marshal(labels)
	if err != nil {
		return ""
	}
	return string(buf)
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		return New(ctx, config.(*Config))
	}))
}
