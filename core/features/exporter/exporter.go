package exporter

import "github.com/VicRen/minidevops/core/features"

type Handler interface {
	features.Feature

	CounterInc(name string, labels map[string]string) error
	CounterAdd(name string, labels map[string]string) error
	GaugeSet(name string, labels map[string]string, value float64) error
	HistogramObserve(name string, labels map[string]string, value float64) error
	SummaryObserve(name string, labels map[string]string, value float64) error
}

func HandlerType() interface{} {
	return (*Handler)(nil)
}
