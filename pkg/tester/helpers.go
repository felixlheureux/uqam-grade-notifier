package tester

import (
	"encoding/json"
	"log"
	"testing"
)

const testcontainersReaper = "testcontainers/ryuk"

// init configure the logger
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func JSONUnmarshal(t *testing.T, data []byte, v interface{}) {
	t.Helper()

	err := json.Unmarshal(data, v)

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

// // AssertEqual compares 2 structs while:
// // * rounding the time-like fields to the nearest second
// // * ignores UpdatedAt & CreatedAt fields
// //
// // This method facilitates comparison when writing/reading structs to Postgres
// func AssertEqual(t *testing.T, expected, actual interface{}, args ...interface{}) {
// 	t.Helper()

// 	opts := []cmp.Option{
// 		cmp.Transformer("T1", func(in time.Time) time.Time {
// 			return in.Round(time.Second)
// 		}),
// 		cmp.Transformer("T2", func(in null.Time) null.Time {
// 			out := null.NewTime(in.Time.Round(time.Second), in.Valid)
// 			return out
// 		}),
// 		cmp.FilterPath(func(p cmp.Path) bool {
// 			for _, ps := range p {
// 				if sf, ok := ps.(cmp.StructField); ok {
// 					if sf.Name() == "CreatedAt" {
// 						return true
// 					}
// 					if sf.Name() == "UpdatedAt" {

// 					}
// 					return true
// 				}
// 			}

// 			return false
// 		}, cmp.Ignore()),
// 	}

// 	diff := cmp.Diff(expected, actual, opts...)

// 	if diff != "" {
// 		var msg string

// 		if len(args) == 1 {
// 			msg = args[0].(string)
// 		}

// 		if len(args) > 1 {
// 			msg = fmt.Sprintf(args[0].(string), args[1:]...)
// 		}

// 		t.Errorf("Assert Not Equal: %s\n%s", msg, diff)
// 	}
// }
