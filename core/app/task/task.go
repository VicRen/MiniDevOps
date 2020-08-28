package task

import (
	"context"
	"fmt"
	"time"

	"github.com/vicren/minidevops/core"
	"github.com/vicren/minidevops/core/common"
	"github.com/vicren/minidevops/core/features/exporter"
)

type Handler struct {
	exporter exporter.Handler
	ticker   *time.Ticker
}

func (h *Handler) Init(config *Config, exporter exporter.Handler) error {
	h.exporter = exporter
	return nil
}

func (h *Handler) Type() interface{} {
	return (*Handler)(nil)
}

func (h *Handler) Start() error {
	h.ticker = time.NewTicker(5 * time.Second)
	count := 0.
	go func() {
		for {
			select {
			case <-h.ticker.C:
				if count > 50 {
					count = 1
				}
				if err := h.exporter.GaugeSet("testing_by_mini_devops_gauge", map[string]string{"testing1": "testing1", "testing2": "string2"}, count); err != nil {
					panic(err)
				}
				count += 5
				fmt.Println("---->task tick", count)
			}
		}
	}()
	return nil
}

func (h *Handler) Close() error {
	h.ticker.Stop()
	return nil
}

func init() {
	common.Must(common.RegisterConfig((*Config)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		h := &Handler{}
		if err := core.RequireFeatures(ctx, func(exporter exporter.Handler) error {
			return h.Init(config.(*Config), exporter)
		}); err != nil {
			return nil, err
		}
		return h, nil
	}))
}
