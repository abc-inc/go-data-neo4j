# Roland

_Roland_ is a **R**epository for **O**bjects and **L**ibrary for **A**ccessing **N**eo4j **D**atabases.
It is heavily inspired by [Spring Data Neo4j][] and aims to be a lightweight alternative to [GoGM][].

## Features

- unified authentication schema for Basic, Bearer and Kerberos authentication
- `Session` with implicit and explicit transaction management
- `Template` for querying `Records` with provided transaction and error handling (similar to [Neo4jTemplate][])
- `Mapper` for mapping each `Record` to a concrete entity or primitive type
- fetching `Metadata` about nodes, relationships and their properties as well as functions and procedures
- make use of [APOC][], if installed, and fallback implementation
- model for accessing execution plans (`EXPLAIN` and `PROFILE`) as well as query statistics

## Roadmap

- implement remaining applicable methods from [Neo4jTemplate]
- implement [Spring Data Neo4j][] inspired `Repository`

## Why _Roland_?

The Neo4j ecosystem has a couple of references to _The Matrix_ franchise e.g., Cypher, Apoc, etc.

Roland is a fictional character in _The Matrix_ franchise.
He was the hard-boiled veteran captain of the Mjolnir during the sixth Matrix Resistance.
Moreover, Mjölnir is the hammer of the thunder god Thor in Norse mythology.
Coincidentally, Neo4j, Inc is a Swedish company and so the circle is complete.

Therefore, _Roland_ is the ideal name for a Neo4j data access library.

[APOC]: https://neo4j.com/developer/neo4j-apoc/
[GoGM]: https://github.com/mindstand/gogm/
[neo4j-go-driver]: github.com/neo4j/neo4j-go-driver/
[Neo4jTemplate]: https://docs.spring.io/spring-data/neo4j/docs/current/api/org/springframework/data/neo4j/core/Neo4jTemplate.html
[Spring Data Neo4j]: https://docs.spring.io/spring-data/neo4j/docs/current/reference/html/
