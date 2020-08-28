package features

import (
	"reflect"

	"github.com/vicren/minidevops/core/common"
)

// Feature is the interface for Galaxy features. All features must implement this interface.
// All existing features have an implementation in app directory. These features can be replaced by third-party ones.
type Feature interface {
	common.HasType
	common.Runnable
}

// GetFeature ...
func GetFeature(allFeatures []Feature, t reflect.Type) Feature {
	for _, f := range allFeatures {
		if reflect.TypeOf(f.Type()) == t {
			return f
		}
	}
	return nil
}
