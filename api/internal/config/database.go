package config

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
)

var (
	Keyspace                = "nihongowa"
	Session  *gocql.Session = nil
)

// Init initializes the database connection
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
	var cluster *gocql.ClusterConfig

	if os.Getenv("ENVIRONMENT") == "prod" {
		cluster := gocql.NewCluster("cassandra.us-east-1.amazonaws.com:9142")
		var auth sigv4.AwsAuthenticator = sigv4.NewAwsAuthenticator()

		auth.Region = os.Getenv("AWS_REGION")
		auth.AccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
		auth.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		cluster.Authenticator = auth

		cluster.Consistency = gocql.LocalQuorum
		cluster.DisableInitialHostLookup = true
	} else {
		cluster_name := os.Getenv("CASSANDRA_CLUSTER_NAME")
		cluster = gocql.NewCluster(cluster_name)

		cluster.Keyspace =
			"nihongowa"
	}

	return cluster
}
