package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"nihongowa/internal/config"
	"nihongowa/internal/handlers"

	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	// Check if the application is running in a Docker environment
	if os.Getenv("ENVIRONMENT") != "Docker" {
		// Attempt to load .env file if not running in Docker
		err := godotenv.Load("../../.env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	connectToCassandra(0)
	config.OpenAIInit()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/messages/:conversation_id", handlers.GetMessagesFromConversation)
	e.POST("/messages/:conversation_id", handlers.PostMessageToConversation)
	e.POST("/conversations", handlers.PostConversation)
	e.GET("/conversations", handlers.GetLastConversations)

	e.Logger.Fatal(e.Start(":1323"))
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

	config.Init(session)
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
