package exporter

import (
	"context"

	"github.com/VicRen/minidevops/core/features/exporter"
	"github.com/prometheus/client_golang/prometheus"
)

type Handler struct {
	counterMap    map[string]prometheus.Counter
	gaugeMap      map[string]prometheus.Gauge
	histogramMap  map[string]prometheus.Histogram
	summaryMap    map[string]prometheus.Summary
	summaryVecMap map[string]*prometheus.SummaryVec
}

func New(ctx context.Context, config *Config) (*Handler, error) {
	h := &Handler{
		counterMap: make(map[string]prometheus.Counter),
	}
	return h, nil
}

func (h *Handler) Type() interface{} {
	return exporter.HandlerType()
}

func (h *Handler) Start() error {
	return nil
}

func (h *Handler) Close() error {
	return nil
}

func (h *Handler) CounterInc(name string, labels map[string]string) error {
	return nil
}

func (h *Handler) CounterAdd(name string, labels map[string]string) error {
	return nil
}

func (h *Handler) GaugeSet(name string, labels map[string]string, value float64) error {
	return nil
}

func (h *Handler) HistogramObserve(name string, labels map[string]string, value float64) error {
	return nil
}

func (h *Handler) SummaryObserve(name string, labels map[string]string, value float64) error {
	return nil
}

func init() {
}
