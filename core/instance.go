package core

import (
	"reflect"
	"sync"

	"github.com/VicRen/minidevops/core/common"
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
