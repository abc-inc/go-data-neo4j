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
	"sort"

	"github.com/abc-inc/roland/graph"
	"golang.org/x/exp/slices"
)

// Metadata holds information about nodes, relationships, properties, functions
// and procedures.
type Metadata struct {
	Nodes []Node
	Rels  []Relationship
	Funcs []Func
	Procs []Func
	Props []string
}

// FetchMetadata retrieves schema information like labels, relationships,
// properties, functions and procedures.
func FetchMetadata(c *graph.Conn) (m Metadata, err error) {
	if c.DBName == graph.SystemDB {
		return
	}

	const cypF = "CALL dbms.functions() YIELD name, signature RETURN name, signature ORDER BY toLower(name)"
	m.Funcs, err = listFuncs(c, cypF)
	if err != nil {
		return m, err
	}

	const cypP = "CALL dbms.procedures() YIELD name, signature RETURN name, signature ORDER BY toLower(name)"
	m.Procs, err = listFuncs(c, cypP)
	if err != nil {
		return m, err
	}

	idx := sort.Search(len(m.Procs), func(i int) bool {
		return m.Procs[i].Name >= "apoc.meta.schema"
	})

	hasApoc := idx < len(m.Procs) && m.Procs[idx].Name == "apoc.meta.schema"
	if !hasApoc {
		return m, fetchSimple(c, &m)
	}
	return m, apocMetaSchema(c, &m)
}

// apocMetaSchema examines a subset of the graph to provide meta information.
func apocMetaSchema(c *graph.Conn, m *Metadata) error {
	res, err := c.Session().Run("CALL apoc.meta.schema", nil)
	if err != nil {
		return err
	}
	rec, err := res.Single()
	if err != nil {
		return err
	}

	var present any
	allProps := make(map[string]any, 0)
	kvs := rec.Values[0].(map[string]any)
	for k, v := range kvs {
		kv := v.(map[string]any)

		if kv["type"] == "node" {
			props := make(map[string]NodeProperty)
			for n, pm := range kv["properties"].(map[string]any) {
				pi := pm.(map[string]any)
				props[n] = NodeProperty{
					Indexed:   pi["indexed"].(bool),
					Unique:    pi["unique"].(bool),
					Existence: pi["existence"].(bool),
					Type:      pi["type"].(string),
				}
				allProps[n] = present
			}

			m.Nodes = append(m.Nodes, Node{
				Count:         kv["count"].(int64),
				Relationships: nil,
				Type:          kv["type"].(string),
				Properties:    props,
				Labels:        []string{k},
			})
		} else {
			props := make(map[string]RelProperty)
			for n, pm := range kv["properties"].(map[string]any) {
				ri := pm.(map[string]any)
				props[n] = RelProperty{
					Array:     ri["array"].(bool),
					Indexed:   ri["indexed"].(bool),
					Existence: ri["existence"].(bool),
					Type:      ri["type"].(string),
				}
				allProps[n] = present
			}

			m.Rels = append(m.Rels, Relationship{
				Count:      kv["count"].(int64),
				Type:       k,
				Properties: props,
			})
		}
	}

	for k := range allProps {
		m.Props = append(m.Props, k)
	}
	sort.Strings(m.Props)

	slices.SortFunc(m.Nodes, func(a, b Node) bool {
		return slices.Compare(a.Labels, b.Labels) < 0
	})
	slices.SortFunc(m.Rels, func(a, b Relationship) bool {
		return a.Type < b.Type
	})

	return nil
}

// fetchSimple loads labels, relationship types and properties.
func fetchSimple(c *graph.Conn, m *Metadata) error {
	const cypLs = "CALL db.labels() YIELD label RETURN label ORDER BY label"
	res, err := c.Session().Run(cypLs, nil)
	if err != nil {
		return err
	}
	for res.Next() {
		l := res.Record().Values[0].(string)
		m.Nodes = append(m.Nodes, Node{Labels: []string{l}})
	}

	cypRels := "CALL db.relationshipTypes() YIELD relationshipType RETURN relationshipType ORDER BY relationshipType"
	res, _ = c.Session().Run(cypRels, nil)
	for res.Next() {
		t := res.Record().Values[0].(string)
		m.Rels = append(m.Rels, Relationship{Type: t})
	}

	const cypProps = "CALL db.propertyKeys() YIELD propertyKey RETURN propertyKey ORDER BY propertyKey"
	res, _ = c.Session().Run(cypProps, nil)
	for res.Next() {
		p := res.Record().Values[0].(string)
		m.Props = append(m.Props, p)
	}

	return nil
}

// listFuncs runs a Cypher statement to list certain functions or procedures.
func listFuncs(c *graph.Conn, cyp string) (funcs []Func, err error) {
	res, err := c.Session().Run(cyp, nil)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		f := Func{
			Name: res.Record().Values[0].(string),
			Sig:  res.Record().Values[1].(string),
		}
		funcs = append(funcs, f)
	}
	return
}
