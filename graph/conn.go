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

import (
	"errors"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const SystemDB = "system"

// defConn holds the default database connection.
var defConn *Conn

// Conn represents a database connection, which can open multiple Sessions.
type Conn struct {
	Driver neo4j.Driver
	user   string
	auth   neo4j.AuthToken
	DBName string
	Tx     neo4j.Transaction
	Params map[string]any
}

// IsConnected returns whether the database connection is established.
func IsConnected() bool {
	return defConn != nil && defConn.DBName != ""
}

// GetConn returns the default connection, regardless of its connection state.
// It panics if there is no connection.
func GetConn() *Conn {
	if defConn == nil {
		panic("Not connected to Neo4j")
	}
	return defConn
}

// NewConn creates a new Neo4j Driver and returns the new Conn.
func NewConn(addr string, user string, auth neo4j.AuthToken, dbName string,
	opts ...func(config *neo4j.Config)) (*Conn, error) {

	d, err := neo4j.NewDriver(addr, auth, opts...)
	if err != nil {
		return nil, err
	}

	defConn = nil
	conn := &Conn{
		Driver: d,
		user:   user,
		auth:   auth,
		DBName: dbName,
		Params: make(map[string]any),
	}

	err = conn.UseDB(dbName)
	if err == nil {
		defConn = conn
	}
	return conn, err
}

// Close the driver and all underlying connections.
func (c *Conn) Close() (err error) {
	if c.Driver != nil {
		if err = c.Driver.Close(); err == nil {
			c.Driver, c.Tx = nil, nil
			c.Params = make(map[string]any)
			c.DBName = ""
		}
	}
	return err
}

// Session creates a new Session.
func (c *Conn) Session() neo4j.Session {
	cfg := neo4j.SessionConfig{DatabaseName: c.DBName}
	return c.Driver.NewSession(cfg)
}

// GetTransaction returns the current Transaction or creates a new one.
func (c *Conn) GetTransaction() (tx neo4j.Transaction, created bool, err error) {
	if c.Tx == nil {
		c.Tx, err = c.Session().BeginTransaction()
		created = true
	}
	return c.Tx, created, err
}

// Commit commits the current Transaction.
// If there is no active Transaction, false is returned.
func (c *Conn) Commit() (done bool, err error) {
	if c.Tx != nil {
		err = c.Tx.Commit()
		c.Tx, done = nil, err != nil
	}
	return
}

// Rollback rolls back the current Transaction.
// If there is no active Transaction, false is returned.
func (c *Conn) Rollback() (done bool, err error) {
	if c.Tx != nil {
		err = c.Tx.Rollback()
		c.Tx, done = nil, err != nil
	}
	return
}

// UseDB permanently changes the database.
func (c *Conn) UseDB(dbName string) (err error) {
	if _, err = c.Rollback(); err != nil {
		return err
	}

	currDBName := c.DBName
	c.DBName = dbName
	_, err = c.Session().ReadTransaction(func(tx neo4j.Transaction) (any, error) {
		return tx.Run("CALL db.ping()", nil)
	})
	var nerr *neo4j.Neo4jError
	if err != nil && errors.As(err, &nerr) {
		if nerr.Title() == "CredentialsExpired" && dbName == SystemDB {
			return nil
		}
	} else if err != nil {
		c.DBName = currDBName
	}
	return err
}

// Username returns the username used to connect to the database.
// If an error occurs, an empty string is returned.
func (c *Conn) Username() string {
	if c.user == "" {
		u, _ := NewTemplate[string](c).QuerySingle(
			"CALL dbms.showCurrentUser()", nil, NewSingleValueMapper[string](0))
		c.user = u
	}
	return c.user
}
