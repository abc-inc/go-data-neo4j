// Copyright 2022 The Roland authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package graph

import "github.com/neo4j/neo4j-go-driver/v4/neo4j/db"

const (
	ID      = "@id"
	Label   = "@label"
	Labels  = "@labels"
	Type    = "@type"
	StartID = "@startId"
	EndID   = "@endId"
)

// Result is an alias for Record with additional convenience methods.
type Result db.Record

// Add appends a new key-value pair to this Result.
func (r *Result) Add(key string, val any) {
	r.Keys = append(r.Keys, key)
	r.Values = append(r.Values, val)
}

// Index returns the position of the given key.
func (r *Result) Index(key string) int {
	for i, k := range r.Keys {
		if k == key {
			return i
		}
	}
	return -1
}

// Value returns the value for the given key.
func (r *Result) Value(key string) (any, bool) {
	idx := r.Index(key)
	if idx < 0 {
		return nil, false
	}
	return r.Values[idx], true
}
