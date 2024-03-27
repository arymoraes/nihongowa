package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"nihongowa/internal/config"
	"nihongowa/internal/handlers"

	"github.com/aws/aws-sigv4-auth-cassandra-gocql-driver-plugin/sigv4"
	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
)

func main() {
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
	var cluster *gocql.ClusterConfig

	if os.Getenv("ENVIRONMENT") == "prod" {
		cluster := gocql.NewCluster("cassandra.us-east-1.amazonaws.com:9142")
		var auth sigv4.AwsAuthenticator = sigv4.NewAwsAuthenticator()

		auth.Region = os.Getenv("AWS_REGION")
		auth.AccessKeyId = os.Getenv("AWS_ACCESS_KEY_ID")
		auth.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")

		cluster.Authenticator = auth

		// cluster.SslOpts = &gocql.SslOptions{
		// 	CaPath: "/Users/user1/.cassandra/AmazonRootCA1.pem",
		// }
		cluster.Consistency = gocql.LocalQuorum
		cluster.DisableInitialHostLookup = true
	}

	if os.Getenv("ENVIRONMENT") != "prod" {
		// "localhost" if dev, "cassandra" if docker
		cluster_name := os.Getenv("CASSANDRA_CLUSTER_NAME")
		cluster = gocql.NewCluster(cluster_name)

		cluster.Keyspace =
			"nihongowa"
	}

	if retryAttempt > 5 {
		panic("Failed to connect to Cassandra")
	}

	session, err := cluster.CreateSession()

	if err != nil {
		time.Sleep(20 * time.Second)
		fmt.Println("Failed to connect to Cassandra, retrying...", err)
		connectToCassandra(retryAttempt + 1)
		return
	}

	config.Init(session)
}
