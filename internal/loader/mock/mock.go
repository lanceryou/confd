package mock

import (
	"fmt"
	"sync"
	"time"

	"github.com/lanceryou/confd/loader"
)

// NewMockLoader new mock loader
func NewMockLoader() *MockLoader {
	return &MockLoader{
		kv:      make(map[string][]byte),
		watchKV: make(map[string]loader.OpType),
	}
}

// MockLoader
type MockLoader struct {
	sync.Mutex
	kv map[string][]byte

	watchKV map[string]loader.OpType
}

// Set write kv
func (m *MockLoader) Set(k string, v []byte) {
	m.Lock()
	defer m.Unlock()

	_, ok := m.kv[k]
	if !ok {
		m.watchKV[k] = loader.Create
	} else if len(v) == 0 {
		m.watchKV[k] = loader.Delete
	} else {
		m.watchKV[k] = loader.Update
	}
	m.kv[k] = v

}

// Load load data by path
func (m *MockLoader) Load(path string) (data []byte, err error) {
	m.Lock()
	defer m.Unlock()

	data, ok := m.kv[path]
	if !ok {
		return nil, fmt.Errorf("load path:%v failed", path)
	}

	return
}

// Watch watch path change, just once
func (m *MockLoader) Watch(path string) (ret *loader.WatchResult, err error) {
	ret = &loader.WatchResult{}
	if op, ok := m.watchKV[path]; ok {
		delete(m.watchKV, path)
		ret.Action = op
		ret.Result = m.kv[path]
		return
	}

	tk := time.NewTicker(time.Second)
	for {
		select {
		case <-tk.C:
			if op, ok := m.watchKV[path]; ok {
				delete(m.watchKV, path)
				ret.Action = op
				ret.Result = m.kv[path]
				return
			}
		}
	}
}

func (m *MockLoader) String() string {
	return "mock"
}
