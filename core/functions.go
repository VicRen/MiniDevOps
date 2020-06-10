package core

import (
	"context"

	"github.com/VicRen/minidevops/core/common"
)

// CreateObject creates a new object based on the given Galaxy instance and config. The Galaxy instance may be nil.
func CreateObject(v Server, config interface{}) (interface{}, error) {
	ctx := context.Background()
	if v != nil {
		ctx = context.WithValue(ctx, instanceKey, v)
	}
	return common.CreateObject(ctx, config)
}
