package main

import (
	"net/http"

	"github.com/charmbracelet/log"

	"nihongowa/internal/config"
	"nihongowa/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.Bootstrap()

	e := echo.New()
	e.Use(middleware.Logger())

	log.Info("Starting server")

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/messages/:conversation_id", handlers.GetMessagesFromConversation)
	e.POST("/messages/:conversation_id", handlers.PostMessageToConversation)
	e.POST("/conversations", handlers.PostConversation)
	e.GET("/conversations", handlers.GetLastConversations)

	e.Logger.Fatal(e.Start(":1323"))
}
