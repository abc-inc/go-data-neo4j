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

package meta

import (
	"strings"
)

// Func represents a function or procedure with its name and signature.
type Func struct {
	Name     string
	Sig      string
	RetItems []Func
}

// String returns the function name.
func (f Func) String() string {
	return f.Name
}

// Schema holds metadata about the database.
type Schema struct {
	Labels   []string
	RelTypes []string
	PropKeys []string
	Funcs    []Func
	Procs    []Func
}

// Node describes the metamodel of a kind of nodes.
type Node struct {
	Count         int64                   `json:"count" yaml:"count" view:"Count"`
	Relationships map[string]NodeRelInfo  `json:"relationships" yaml:"relationships" view:"Relationships"`
	Type          string                  `json:"type" yaml:"type" view:"Type"`
	Properties    map[string]NodeProperty `json:"properties" yaml:"properties" view:"Properties"`
	Labels        []string                `json:"labels" yaml:"labels" view:"Labels"`
}

// String returns the labels of the Node.
func (n Node) String() string {
	return ":" + strings.Join(n.Labels, ":")
}

// Relationship describes the metamodel of one kind of relationships.
type Relationship struct {
	Count      int64                  `json:"count" yaml:"count" view:"Count"`
	Type       string                 `json:"type" yaml:"type" view:"Type"`
	Properties map[string]RelProperty `json:"properties" yaml:"properties" view:"Properties"`
}

// String returns the relationship type.
func (r Relationship) String() string {
	return ":" + r.Type
}

// NodeRelInfo describes a Relationship type of one kind of nodes.
type NodeRelInfo struct {
	Count      int                    `json:"count" yaml:"count" view:"Count"`
	Properties map[string]RelProperty `json:"properties" yaml:"properties" view:"Properties"`
	Direction  string                 `json:"direction" yaml:"direction" view:"Direction"`
	Labels     []string               `json:"labels" yaml:"labels" view:"Labels"`
}

// NodeProperty represents a single property of a Node.
type NodeProperty struct {
	Indexed   bool   `json:"indexed" yaml:"indexed" view:"Indexed"`
	Unique    bool   `json:"unique" yaml:"unique" view:"Unique"`
	Existence bool   `json:"existence" yaml:"existence" view:"Existence"`
	Type      string `json:"type" yaml:"type" view:"Type"`
}

// RelProperty represents a single property of a Relationship.
type RelProperty struct {
	Array     bool   `json:"array" yaml:"array" view:"Array"`
	Indexed   bool   `json:"indexed" yaml:"indexed" view:"Indexed"`
	Existence bool   `json:"existence" yaml:"existence" view:"Existence"`
	Type      string `json:"type" yaml:"type" view:"Type"`
}
