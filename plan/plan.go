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

package plan

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// StmtType is an alias for StatementType.
// It is used to allow custom string representation and marshalling.
type StmtType neo4j.StatementType

// String describes the type of statement.
func (s StmtType) String() string {
	return []string{"UNKNOWN", "READ_ONLY", "READ_WRITE", "WRITE_ONLY", "SCHEMA_WRITE"}[s]
}

// MarshalJSON serializes the statement type as JSON value.
func (s StmtType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + s.String() + "\""), nil
}

// MarshalYAML serializes the statement type as YAML value.
func (s StmtType) MarshalYAML() (any, error) {
	return s.String(), nil
}

// Stats holds the statics of an execution plan.
type Stats struct {
	Plan      string   `json:"plan" yaml:"plan" view:"Plan"`
	Statement StmtType `json:"queryType" yaml:"queryType" view:"Statement"`
	Version   string   `json:"version" yaml:"version" view:"Version"`
	Planner   string   `json:"planner" yaml:"planner" view:"Planner"`
	Runtime   string   `json:"runtime" yaml:"runtime" view:"Runtime"`
	Time      int64    `json:"time" yaml:"time" view:"Time"`
	DBHits    int64    `json:"dbHits" yaml:"dbHits" view:"DB Hits"`
	Rows      int64    `json:"rows" yaml:"rows" view:"Rows"`
	Memory    int64    `json:"memory" yaml:"memory" view:"Memory (Bytes)"`
}

// String returns a string representation of the statistics.
func (p Stats) String() string {
	return fmt.Sprintf("%s (%d ms, %d rows, %d DB hits)", p.Plan, p.Time, p.Rows, p.DBHits)
}

// Op is a single operation in an execution plan.
type Op struct {
	Op          string `json:"operatorType" yaml:"operatorType" view:"Operator"`
	Details     string `json:"details,omitempty" yaml:"details,omitempty" view:"Details,omitempty"`
	RowsEst     int64  `json:"estimatedRows" yaml:"estimatedRows" view:"Estimated Rows"`
	Rows        int64  `json:"rows" yaml:"rows" view:"Rows"`
	DBHits      int64  `json:"dbHits" yaml:"dbHits" view:"DB Hits"`
	Memory      int64  `json:"memory" yaml:"memory" view:"Memory (Bytes)"`
	CacheHits   int64  `json:"pageCacheHits" yaml:"pageCacheHits" view:"Cache Hits"`
	CacheMisses int64  `json:"pageCacheMisses" yaml:"pageCacheMisses" view:"Cache Misses"`
	Order       string `json:"order,omitempty" yaml:"order,omitempty" view:"Ordered by,omitempty"`
	Children    []*Op  `json:"children,omitempty" yaml:"children,omitempty" view:"-"`
}

// String returns a string representation of the operation.
func (p Op) String() string {
	return fmt.Sprintf("%s (%d %s, %d DB %s)", p.Op,
		p.Rows, plural(p.Rows, "row", "rows"),
		p.DBHits, plural(p.DBHits, "hit", "hits"))
}

// plural returns singular if n is 1 or plural otherwise.
func plural(n int64, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}
