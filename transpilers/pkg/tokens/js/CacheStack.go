package js

import (
	"sync"

	"./values"
)

var SERIAL = false

type CacheStack struct {
	// all values can be saved here
	genInst  map[interface{}]*values.Instance
	cacheVal map[interface{}]values.Value
	mux      *sync.Mutex
}

func NewCacheStack() *CacheStack {
	return &CacheStack{
		make(map[interface{}]*values.Instance),
		make(map[interface{}]values.Value),
		&sync.Mutex{},
	}
}

func (s *CacheStack) SetGeneratedInstance(ptr interface{}, inst *values.Instance) {
	s.mux.Lock()
	s.genInst[ptr] = inst
	s.mux.Unlock()
}

func (s *CacheStack) GetGeneratedInstance(ptr interface{}) (*values.Instance, bool) {
	s.mux.Lock()
	inst, ok := s.genInst[ptr]
	s.mux.Unlock()

	if SERIAL {
		return inst, ok
	} else if ok {
		if inst == nil {
			return nil, true
		} else {
			return inst.Copy(values.NewCopyCache()).(*values.Instance), true
		}
	} else {
		return nil, false
	}
}

func (s *CacheStack) SetCacheValue(ptr interface{}, val values.Value) {
	s.mux.Lock()
	s.cacheVal[ptr] = val
	s.mux.Unlock()
}

func (s *CacheStack) GetCacheValue(ptr interface{}) (values.Value, bool) {
	s.mux.Lock()
	val, ok := s.cacheVal[ptr]
	s.mux.Unlock()

	if SERIAL {
		return val, ok
	} else if ok {
		if val == nil {
			return nil, true
		} else {
			return val.Copy(values.NewCopyCache()), true
		}
	} else {
		return nil, false
	}
}
