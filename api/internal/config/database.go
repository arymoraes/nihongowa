package config

import "github.com/gocql/gocql"

var (
	Keyspace                = "nihongowa"
	Session  *gocql.Session = nil
)

// Init initializes the database connection
func Init(session *gocql.Session) {
	Session = session
	createSchema(session)
}
