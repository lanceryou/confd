package confd

import (
	"testing"
)

func TestParseSchema(t *testing.T) {
	ts := []struct {
		schema       string
		expectSchema *Schema
		err          error
	}{
		{
			schema:       "./config/service.1.yml",
			expectSchema: &Schema{source: "file", format: "yml", key: "./config/service.1.yml"},
		},
		{
			schema: "",
			err:    SchemaErr,
		},
		{
			schema: "etcd:yml",
			err:    SchemaErr,
		},
		{
			schema:       "etcd:yml:service",
			expectSchema: &Schema{source: "etcd", format: "yml", key: "service"},
		},
	}

	for _, s := range ts {
		sche, err := ParseSchema(s.schema)
		if err != s.err {
			t.Errorf("parse schema expect err %v, but err %v", s.err, err)
		}

		if sche != nil && (sche.format != s.expectSchema.format ||
			sche.source != s.expectSchema.source ||
			sche.key != s.expectSchema.key) {
			t.Errorf("expect schema %v, but schema is %v", sche, s.expectSchema)
		}
	}
}
