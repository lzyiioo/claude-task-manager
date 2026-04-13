// internal/core/json.go
package core

import "encoding/json"

// mustMarshal panics if JSON marshaling fails
func mustMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}
