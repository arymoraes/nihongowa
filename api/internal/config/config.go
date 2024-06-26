package config

import (
	"os"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
)

var BasePath string = ""

func Bootstrap() {
	loadEnv()
	loadBasePath()
	connectToCassandra(0)
	openAIInit()
}

func loadEnv() {
	log.Info("Loading environment variables")
	if os.Getenv("ENVIRONMENT") != "prod" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}
}

func createSchema(session *gocql.Session) {
	log.Info("Creating schema")
	createKeyspace(session)

	log.Info("Creating tables")

	// Create the table 'conversations'
	createTableCql := `CREATE TABLE IF NOT EXISTS nihongowa.conversations (
        id UUID PRIMARY KEY,
				assistant_name VARCHAR,
				thread_id VARCHAR,
				assistant_id VARCHAR,
				run_id VARCHAR,
				scenario TEXT
    );`

	// Execute the CQL to create the table
	if err := session.Query(createTableCql).Exec(); err != nil {
		log.Error("Failed to create table 'conversations': %v", err)
	}

	createMessagesTableCql := `CREATE TABLE IF NOT EXISTS nihongowa.messages (
				id UUID,
				content TEXT,
				translation TEXT,
				romanji TEXT,
				user_message_translated TEXT,
				is_ai BOOLEAN,
				word_by_word_translation LIST<TEXT>,
				conversation_id UUID,
				created_at TIMESTAMP,
				updated_at TIMESTAMP,
				PRIMARY KEY (conversation_id, id)
		);`

	if err := session.Query(createMessagesTableCql).Exec(); err != nil {
		log.Error("Failed to create table 'messages': %v", err)
	}
}

func createKeyspace(session *gocql.Session) {
	cql := `CREATE KEYSPACE IF NOT EXISTS nihongowa WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };`

	if err := session.Query(cql).Exec(); err != nil {
		log.Fatalf("Failed to create keyspace: %v", err)
	}

	log.Info("Keyspace created")
}

func loadBasePath() {
	if os.Getenv("BASE_PATH") != "" {
		BasePath = os.Getenv("BASE_PATH") + "/"
	}
}
