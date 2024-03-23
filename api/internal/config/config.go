package config

import (
	"log"

	"github.com/gocql/gocql"
)

func createSchema(session *gocql.Session) {
	createKeyspace(session)

	// Create the custom type 'message'
	createTypeCql := `CREATE TYPE IF NOT EXISTS nihongowa.message (
        content TEXT,
        translation TEXT,
        wordByWordTranslation LIST<TEXT>,
        createdAt TIMESTAMP,
        updatedAt TIMESTAMP
    );`

	// Execute the CQL to create the type
	if err := session.Query(createTypeCql).Exec(); err != nil {
		log.Fatalf("Failed to create custom type 'message': %v", err)
	}

	// Create the table 'conversations'
	createTableCql := `CREATE TABLE IF NOT EXISTS nihongowa.conversations (
        id UUID PRIMARY KEY,
        messages LIST<FROZEN<message>>,
				ThreadID VARCHAR,
				AssistantID VARCHAR,
				Scenario TEXT
    );`

	// Execute the CQL to create the table
	if err := session.Query(createTableCql).Exec(); err != nil {
		log.Fatalf("Failed to create table 'conversations': %v", err)
	}
}

func createKeyspace(session *gocql.Session) {
	cql := `CREATE KEYSPACE IF NOT EXISTS nihongowa WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`

	if err := session.Query(cql).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	log.Println("Keyspace created")
}
