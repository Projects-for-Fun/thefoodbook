# How-to

Create migration files with:

```migrate create -ext cypher -dir db/migrations <filename>```

(Need one statement per file.)

<br>

For neo4j 5+:
```
// up:
CREATE CONSTRAINT FOR (a:SchemaMigration) REQUIRE a.version IS UNIQUE

// down:
DROP CONSTRAINT ON (a:SchemaMigration) ASSERT a.version IS UNIQUE
```

