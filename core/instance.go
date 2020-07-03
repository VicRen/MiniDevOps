package core

import (
	"errors"
	"log"
	"reflect"
	"sync"

	"github.com/VicRen/minidevops/core/common"
	"github.com/VicRen/minidevops/core/common/serial"
	"github.com/VicRen/minidevops/core/features"
)

type Server interface {
	common.Runnable
}

func ServerType() interface{} {
	return (*Instance)(nil)
}

type Instance struct {
	access             sync.Mutex
	features           []features.Feature
	featureResolutions []resolution
	running            bool
}

func New(config *Config) (*Instance, error) {
	var server = &Instance{}

	for _, feature := range config.Features {
		af, err := serial.UnmarshalAny(feature)
		if err != nil {
			return nil, err
		}
		obj, err := CreateObject(server, af)
		if err != nil {
			return nil, err
		}
		if f, ok := obj.(features.Feature); ok {
			if err := server.AddFeature(f); err != nil {
				return nil, err
			}
		}
	}

	if server.featureResolutions != nil {
		return nil, errors.New("not all dependency are resolved")
	}

	return server, nil
}

// AddFeature registers a feature into current Instance.
func (i *Instance) AddFeature(feature features.Feature) error {
	i.features = append(i.features, feature)

	if i.running {
		if err := feature.Start(); err != nil {
			log.Printf("failed to start feature: %v", err)
		}
		return nil
	}

	if i.featureResolutions == nil {
		return nil
	}

	var pendingResolutions []resolution
	for _, r := range i.featureResolutions {
		finished, err := r.resolve(i.features)
		if finished && err != nil {
			return err
		}
		if !finished {
			pendingResolutions = append(pendingResolutions, r)
		}
	}
	if len(pendingResolutions) == 0 {
		i.featureResolutions = nil
	} else if len(pendingResolutions) < len(i.featureResolutions) {
		i.featureResolutions = pendingResolutions
	}

	return nil
}

func (i *Instance) Type() interface{} {
	return ServerType()
}

func (i *Instance) Start() error {
	i.access.Lock()
	defer i.access.Unlock()

	i.running = true

	return nil
}

func (i *Instance) Close() error {
	i.access.Lock()
	defer i.access.Unlock()

	i.running = false

	return nil
}

type resolution struct {
	deps     []reflect.Type
	callback interface{}
}

func (r *resolution) resolve(allFeatures []features.Feature) (bool, error) {
	var fs []features.Feature
	for _, d := range r.deps {
		f := getFeature(allFeatures, d)
		if f == nil {
			return false, nil
		}
		fs = append(fs, f)
	}

	callback := reflect.ValueOf(r.callback)
	var input []reflect.Value
	callbackType := callback.Type()
	for i := 0; i < callbackType.NumIn(); i++ {
		pt := callbackType.In(i)
		for _, f := range fs {
			if reflect.TypeOf(f).AssignableTo(pt) {
				input = append(input, reflect.ValueOf(f))
				break
			}
		}
	}

	if len(input) != callbackType.NumIn() {
		panic("Can't get all input parameters")
	}

	var err error
	ret := callback.Call(input)
	errInterface := reflect.TypeOf((*error)(nil)).Elem()
	for i := len(ret) - 1; i >= 0; i-- {
		if ret[i].Type() == errInterface {
			v := ret[i].Interface()
			if v != nil {
				err = v.(error)
			}
			break
		}
	}

	return true, err
}

func getFeature(allFeatures []features.Feature, t reflect.Type) features.Feature {
	for _, f := range allFeatures {
		if reflect.TypeOf(f.Type()) == t {
			return f
		}
	}
	return nil
}
