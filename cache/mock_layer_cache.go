// Code generated by mockery v1.0.0. DO NOT EDIT.

package cache

import mock "github.com/stretchr/testify/mock"
import types "github.com/aquasecurity/fanal/types"

// MockLayerCache is an autogenerated mock type for the LayerCache type
type MockLayerCache struct {
	mock.Mock
}

type MissingLayersArgs struct {
	Layers         []string
	LayersAnything bool
}

type MissingLayersReturns struct {
	MissingLayerIDs []string
	Err             error
}

type MissingLayersExpectation struct {
	Args    MissingLayersArgs
	Returns MissingLayersReturns
}

func (_m *MockLayerCache) ApplyMissingLayersExpectation(e MissingLayersExpectation) {
	var args []interface{}
	if e.Args.LayersAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.Layers)
	}
	_m.On("MissingLayers", args...).Return(e.Returns.MissingLayerIDs, e.Returns.Err)
}

func (_m *MockLayerCache) ApplyMissingLayersExpectations(expectations []MissingLayersExpectation) {
	for _, e := range expectations {
		_m.ApplyMissingLayersExpectation(e)
	}
}

// MissingLayers provides a mock function with given fields: layers
func (_m *MockLayerCache) MissingLayers(layers []string) ([]string, error) {
	ret := _m.Called(layers)

	var r0 []string
	if rf, ok := ret.Get(0).(func([]string) []string); ok {
		r0 = rf(layers)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(layers)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type PutLayerArgs struct {
	LayerID                     string
	LayerIDAnything             bool
	DecompressedLayerID         string
	DecompressedLayerIDAnything bool
	LayerInfo                   types.LayerInfo
	LayerInfoAnything           bool
}

type PutLayerReturns struct {
	Err error
}

type PutLayerExpectation struct {
	Args    PutLayerArgs
	Returns PutLayerReturns
}

func (_m *MockLayerCache) ApplyPutLayerExpectation(e PutLayerExpectation) {
	var args []interface{}
	if e.Args.LayerIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.LayerID)
	}
	if e.Args.DecompressedLayerIDAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.DecompressedLayerID)
	}
	if e.Args.LayerInfoAnything {
		args = append(args, mock.Anything)
	} else {
		args = append(args, e.Args.LayerInfo)
	}
	_m.On("PutLayer", args...).Return(e.Returns.Err)
}

func (_m *MockLayerCache) ApplyPutLayerExpectations(expectations []PutLayerExpectation) {
	for _, e := range expectations {
		_m.ApplyPutLayerExpectation(e)
	}
}

// PutLayer provides a mock function with given fields: layerID, decompressedLayerID, layerInfo
func (_m *MockLayerCache) PutLayer(layerID string, decompressedLayerID string, layerInfo types.LayerInfo) error {
	ret := _m.Called(layerID, decompressedLayerID, layerInfo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, types.LayerInfo) error); ok {
		r0 = rf(layerID, decompressedLayerID, layerInfo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}