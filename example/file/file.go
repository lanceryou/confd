package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-yaml/yaml"

	"github.com/lanceryou/confd"
)

type NestStruct struct {
	Nick  string  `json:"nick,omitempty" yaml:"nick,omitempty"`
	Score float64 `json:"score,omitempty" yaml:"score,omitempty"`
}

type TestStruct struct {
	Count int64      `json:"count,omitempty" yaml:"count,omitempty"`
	Str   string     `json:"str,omitempty" yaml:"str,omitempty"`
	Nest  NestStruct `json:"nest,omitempty" yaml:"nest,omitempty"`
}

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

func writeFile(path string, content []byte) {
	fileObj, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer fileObj.Close()
	if _, err := fileObj.Write(content); err != nil {
		panic(err)
	}
}

func main() {
	cnd := confd.NewConfd(confd.WithSchema("test1.yml"))

	writeFile("test1.yml", mustYml(
		TestStruct{Count: 1, Str: "test", Nest: NestStruct{Nick: "lancer", Score: 1.2}},
	))
	if err := cnd.LoadConfig(); err != nil {
		panic(err)
	}

	for _, key := range []string{"count", "str", "nest", "nest.score", "nest.nick"} {
		if !cnd.HasValue(key) {
			fmt.Printf("expect has key %v", key)
		}
	}

	for k, v := range map[string]int64{"count": 1} {
		var val int64
		if err := cnd.ReadSection(k, &val); err != nil {
			fmt.Printf("read error %v", err)
		}

		if val != v {
			fmt.Printf("expect val %v, but val %v", v, val)
		}
	}

	for k, v := range map[string]float64{"nest.score": 1.2} {
		var val float64
		if err := cnd.ReadSection(k, &val); err != nil {
			fmt.Printf("read error %v", err)
		}

		if val != v {
			fmt.Printf("expect val %v, but val %v", v, val)
		}
	}

	for k, v := range map[string]string{
		"str":       "test",
		"nest.nick": "lancer",
	} {
		var val string
		if err := cnd.ReadSection(k, &val); err != nil {
			fmt.Printf("read error %v", err)
		}

		if val != v {
			fmt.Printf("expect val %v, but val %v", v, val)
		}
	}

	var testStruct TestStruct
	if err := cnd.Read(&testStruct); err != nil {
		panic(err)
	}
	fmt.Printf("test struct %v\n", testStruct)

	go func() {
		t := time.NewTicker(time.Second)
		for {
			select {
			case <-t.C:
				// change file
				writeFile("test1.yml", mustYml(
					TestStruct{Count: 23, Str: "test123", Nest: NestStruct{Nick: "lancerwith", Score: 1.2}},
				))
			}
		}
	}()

	if err := cnd.WatchConfig(); err != nil {
		fmt.Printf("watch err %v", err)
	}

	fmt.Printf("receive file change.\n")
	// wait watch change
	for k, v := range map[string]int64{"count": 23} {
		var val int64
		if err := cnd.ReadSection(k, &val); err != nil {
			fmt.Printf("read error %v", err)
		}

		if val != v {
			fmt.Printf("expect val %v, but val %v", v, val)
		}
	}

	for k, v := range map[string]float64{"nest.score": 1.2} {
		var val float64
		if err := cnd.ReadSection(k, &val); err != nil {
			fmt.Printf("read error %v", err)
		}

		if val != v {
			fmt.Printf("expect val %v, but val %v", v, val)
		}
	}

	for k, v := range map[string]string{
		"str":       "test123",
		"nest.nick": "lancerwith",
	} {
		var val string
		if err := cnd.ReadSection(k, &val); err != nil {
			fmt.Printf("read error %v", err)
		}

		if val != v {
			fmt.Printf("expect val %v, but val %v", v, val)
		}
	}

	if err := cnd.Read(&testStruct); err != nil {
		panic(err)
	}
	fmt.Printf("test struct %v\n", testStruct)
}
