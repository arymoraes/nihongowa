package config

import (
	"fmt"
	"time"

	"nihongowa/internal/utils"

	"github.com/gocql/gocql"
	"github.com/labstack/gommon/log"
)

var (
	Keyspace                = "nihongowa"
	Session  *gocql.Session = nil
)

func initSession(session *gocql.Session) {
	Session = session
	createSchema(session)
}

func connectToCassandra(retryAttempt int) {
	cluster := configureCassandraCluster()

	if retryAttempt > 5 {
		panic("Failed to connect to Cassandra")
	}

	session, err := cluster.CreateSession()

	if err != nil {
		fmt.Println("Failed to connect to Cassandra, retrying...", err)
		time.Sleep(20 * time.Second)
		connectToCassandra(retryAttempt + 1)
		return
	}

	initSession(session)
}

func configureCassandraCluster() *gocql.ClusterConfig {
	// if os.Getenv("ENVIRONMENT") == "prod" {
	// 	log.Info("Configuring AWS Keyspaces")
	// 	return configureKeyspaces()
	// } else {
	log.Info("Configuring local Cassandra")
	cluster_name := utils.GetEnv("CASSANDRA_CLUSTER_NAME", "localhost")
	cluster := gocql.NewCluster(cluster_name)

	cluster.Keyspace =
		"nihongowa"

	return cluster
	// }
}
