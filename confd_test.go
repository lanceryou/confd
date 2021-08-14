package confd

import (
	"encoding/json"
	"testing"

	"github.com/go-yaml/yaml"

	"github.com/lanceryou/confd/internal/loader/mock"
	"github.com/lanceryou/confd/loader"
)

func mustJson(val interface{}) []byte {
	data, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}

	return data
}

func mustYml(val interface{}) []byte {
	data, err := yaml.Marshal(val)
	if err != nil {
		panic(err)
	}

	return data
}

var (
	ml = mock.NewMockLoader()
)

func registerLoader() {
	loader.RegisterConfigLoader(ml)
}

type NestStruct struct {
	Nick  string  `json:"nick,omitempty" yaml:"nick,omitempty"`
	Score float64 `json:"score,omitempty" yaml:"score,omitempty"`
}

type TestStruct struct {
	Count int64       `json:"count,omitempty" yaml:"count,omitempty"`
	Str   string      `json:"str,omitempty" yaml:"str,omitempty"`
	Nest  *NestStruct `json:"nest,omitempty" yaml:"nest,omitempty"`
}

func TestConfd_LoadConfig(t *testing.T) {
	registerLoader()

	ts := []struct {
		testStruct  TestStruct
		key         string
		schema      string
		fn          func(interface{}) []byte
		hasKeys     []string
		excludeKeys []string
		int64Map    map[string]int64
		strMap      map[string]string
		float64Map  map[string]float64
	}{
		{
			testStruct: TestStruct{Count: 1, Str: "test", Nest: &NestStruct{Nick: "lancer", Score: 1.2}},
			key:        "key",
			schema:     "mock:json:key",
			fn:         mustJson,
			hasKeys:    []string{"count", "str", "nest", "nest.score", "nest.nick"},
			int64Map: map[string]int64{
				"count": 1,
			},
			strMap: map[string]string{
				"str":       "test",
				"nest.nick": "lancer",
			},
			float64Map: map[string]float64{
				"nest.score": 1.2,
			},
		},
		{
			testStruct: TestStruct{Count: 2, Str: "test", Nest: &NestStruct{Nick: "lancer", Score: 1.2}},
			key:        "key",
			schema:     "mock:yaml:key",
			fn:         mustYml,
			hasKeys:    []string{"count", "str", "nest", "nest.score", "nest.nick"},
			int64Map: map[string]int64{
				"count": 2,
			},
			strMap: map[string]string{
				"str":       "test",
				"nest.nick": "lancer",
			},
			float64Map: map[string]float64{
				"nest.score": 1.2,
			},
		},
		{
			key:         "key",
			schema:      "mock:json:key",
			fn:          mustJson,
			excludeKeys: []string{"count", "str", "nest", "nest.score", "nest.nick"},
		},
		{
			key:         "key",
			schema:      "mock:yaml:key",
			fn:          mustYml,
			excludeKeys: []string{"count", "str", "nest", "nest.score", "nest.nick"},
		},
	}

	for _, s := range ts {
		conf := NewConfd(WithSchema(s.schema))
		ml.Set(s.key, s.fn(s.testStruct))
		if err := conf.LoadConfig(); err != nil {
			t.Errorf("some err:%v", err)
		}

		for _, key := range s.hasKeys {
			if !conf.HasValue(key) {
				t.Errorf("expect has key:%v", key)
			}
		}

		for _, key := range s.excludeKeys {
			if conf.HasValue(key) {
				t.Errorf("expect has key:%v", key)
			}
		}

		for k, v := range s.int64Map {
			var val int64
			if err := conf.ReadSection(k, &val); err != nil {
				t.Errorf("read error %v", err)
			}

			if val != v {
				t.Errorf("expect val %v, but val %v", v, val)
			}
		}

		for k, v := range s.float64Map {
			var val float64
			if err := conf.ReadSection(k, &val); err != nil {
				t.Errorf("read error %v", err)
			}

			if val != v {
				t.Errorf("expect val %v, but val %v", v, val)
			}
		}

		for k, v := range s.strMap {
			var val string
			if err := conf.ReadSection(k, &val); err != nil {
				t.Errorf("read error %v", err)
			}

			if val != v {
				t.Errorf("expect val %v, but val %v", v, val)
			}
		}
	}

}

func TestConfd_WatchConfig(t *testing.T) {

}
