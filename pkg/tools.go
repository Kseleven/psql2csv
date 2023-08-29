package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"net/netip"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func valueToString(src any) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return formatTime(v)
	case netip.Prefix:
		return src.(netip.Prefix).String()
	case netip.Addr:
		return src.(netip.Addr).String()
	case bool:
		return strconv.FormatBool(v)
	default:
		if v == nil {
			return ""
		}
		if t := reflect.TypeOf(v); t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			v := reflect.ValueOf(v)
			var buf strings.Builder
			buf.WriteString("{")
			for i := 0; i < v.Len(); i++ {
				if v.Index(i).Elem().CanInterface() && v.Index(i).Elem().Type() == reflect.TypeOf(pgtype.Numeric{}) {
					n, ok := v.Index(i).Elem().Interface().(pgtype.Numeric)
					if ok {
						if n.InfinityModifier == pgtype.Finite {
							r, _ := n.Int64Value()
							buf.WriteString(strconv.FormatInt(r.Int64, 10))
						} else {
							r, _ := n.Float64Value()
							buf.WriteString(strconv.FormatFloat(r.Float64, 'E', -1, 64))
						}
					}
				} else {
					buf.WriteString(v.Index(i).Elem().String())
				}

				if i < v.Len()-1 {
					buf.WriteString(",")
				}
			}
			buf.WriteString("}")
			return buf.String()
		}

		b, _ := json.Marshal(v)
		return string(b)
	}
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

func columnDiff(a []string, b []string) error {
	missingCloumns := make([]string, 0, len(b))
	surplusCloumns := make([]string, 0, len(a))
	aMap := make(map[string]bool, len(a))
	for _, s := range a {
		aMap[s] = false
	}

	diff := false
	bMap := make(map[string]bool, len(b))
	for _, s := range b {
		if _, ok := aMap[s]; !ok {
			missingCloumns = append(missingCloumns, s)
			diff = true
		}
		bMap[s] = false
	}

	for _, s := range a {
		if _, ok := bMap[s]; !ok {
			surplusCloumns = append(surplusCloumns, s)
			diff = true
		}
	}

	if diff {
		return fmt.Errorf("missing columns:%s surplus columns:%s", missingCloumns, surplusCloumns)
	}
	return nil
}
