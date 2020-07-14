package task

import (
	"context"
	"fmt"
	"time"

	"github.com/VicRen/minidevops/core"
	"github.com/VicRen/minidevops/core/common"
	"github.com/VicRen/minidevops/core/features/exporter"
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
	go func() {
		for {
			select {
			case <-h.ticker.C:
				if err := h.exporter.CounterInc("testing_by_mini_devops", map[string]string{"testing1": "testing1", "testing2": "string2"}); err != nil {
					panic(err)
				}
				fmt.Println("---->task tick")
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
