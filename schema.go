package confd

import (
	"errors"
	"strings"
)

var (
	SchemaErr = errors.New("schema error")
)

// Schema
// schema {source}:{format}:{key}
// source file, etcd,...
// format ymal, json,ini...
// key can read config information, eg ./xx.yml...
type Schema struct {
	source string
	format string
	key    string
}

// ParseSchema
func ParseSchema(schema string) (s *Schema, err error) {
	ss := strings.Split(schema, ":")
	if len(ss) == 1 {
		return parsePathSchema(schema)
	}

	if len(ss) != 3 {
		return nil, SchemaErr
	}

	return &Schema{
		source: ss[0],
		format: ss[1],
		key:    ss[2],
	}, nil
}

func parsePathSchema(path string) (s *Schema, err error) {
	id := strings.LastIndex(path, ".")
	if id == -1 {
		return nil, SchemaErr
	}

	return &Schema{
		source: "file",
		format: path[id+1:],
		key:    path,
	}, nil

}
