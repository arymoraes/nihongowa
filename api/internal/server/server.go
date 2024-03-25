package main

import (
	"net/http"
	"os"

	"nihongowa/internal/config"
	"nihongowa/internal/handlers"

	"github.com/gocql/gocql"
	"github.com/labstack/echo/v4"
)

func main() {
	cluster_name := "localhost"

	if os.Getenv("ENVIRONMENT") != "Docker" {
		cluster_name = "cassandra"
	}

	cluster := gocql.NewCluster(cluster_name)
	cluster.Keyspace = "nihongowa"
	session, err := cluster.CreateSession()

	if err != nil {
		panic(err)
	}

	config.Init(session)
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
