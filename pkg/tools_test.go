package pkg

import (
	"errors"
	"net/netip"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestColumnDiff(t *testing.T) {
	var testDatas = []struct {
		name             string
		column1, column2 []string
		expect           error
	}{
		{
			name:    "same",
			column1: []string{"name", "address", "age"},
			column2: []string{"name", "address", "age"},
			expect:  nil,
		},
		{
			name:    "same1",
			column1: []string{"id", "create_time", "name", "parent_id", "depth", "path_ids", "read_only", "index", "brief_name", "unique_code", "comment"},
			column2: []string{"id", "create_time", "name", "brief_name", "unique_code", "comment", "parent_id", "depth", "path_ids", "read_only", "index"},
			expect:  nil,
		},
		{
			name:    "missing column",
			column1: []string{"name", "address", "age"},
			column2: []string{"name", "address", "age", "country"},
			expect:  errors.New("missing columns:[country] surplus columns:[]"),
		},
		{
			name:    "surplus column",
			column1: []string{"name", "address", "age"},
			column2: []string{"name", "age"},
			expect:  errors.New("missing columns:[] surplus columns:[address]"),
		},
		{
			name:    "mix column",
			column1: []string{"name", "address", "age"},
			column2: []string{"name", "age", "country"},
			expect:  errors.New("missing columns:[country] surplus columns:[address]"),
		},
		{
			name:    "multiple columns",
			column1: []string{"name", "address", "sex", "age"},
			column2: []string{"name", "age", "country", "city"},
			expect:  errors.New("missing columns:[country city] surplus columns:[address sex]"),
		},
	}

	for _, data := range testDatas {
		t.Run(data.name, func(t *testing.T) {
			assert.Equal(t, columnDiff(data.column1, data.column2), data.expect)
		})
	}
}

func TestValueToString(t *testing.T) {
	var testDatas = []struct {
		name   string
		value  any
		expect any
	}{
		{name: "string", value: "hello", expect: "hello"},
		{name: "byte", value: []byte("hello"), expect: "hello"},
		{name: "int", value: 12, expect: "12"},
		{name: "int32", value: int32(12), expect: "12"},
		{name: "int64", value: int64(12), expect: "12"},
		{name: "uint32", value: uint32(12), expect: "12"},
		{name: "uint64", value: uint64(12), expect: "12"},
		{name: "float32", value: float32(12), expect: "12"},
		{name: "float64", value: 12.2, expect: "12.2"},
		{name: "time.Time", value: func() time.Time {
			return time.Date(2023, 5, 25, 10, 36, 7, 0, time.Local)
		}(), expect: "2023-05-25T10:36:07+08:00"},
		{name: "netip.Prefix", value: netip.MustParsePrefix("10.0.0.0/24"), expect: "10.0.0.0/24"},
		{name: "netip.Addr", value: netip.MustParseAddr("10.0.0.1"), expect: "10.0.0.1"},
		{name: "bool", value: false, expect: "false"},
	}

	for _, data := range testDatas {
		t.Run(data.name, func(t *testing.T) {
			assert.Equal(t, valueToString(data.value), data.expect)
		})
	}
}
