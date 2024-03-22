package handlers

import (
	"net/http"

	"nihongowa/internal/models"
	"nihongowa/internal/services"

	"github.com/labstack/echo/v4"
)

// e.GET("/messages/:conversation_id", getMessagesFromConversation)
func GetMessagesFromConversation(c echo.Context) error {
	id := c.Param("conversation_id")

	messages := services.GetMessagesFromConversation(id)

	return c.JSON(http.StatusOK, messages)
}

// e.POST("/messages/:conversation_id", postMessageToConversation)
func PostMessageToConversation(c echo.Context) error {
	id := c.Param("conversation_id")

	content := c.FormValue("message")

	message := models.Message{
		Content: content,
	}

	services.PostMessageToConversation(id, message)

	return c.String(http.StatusOK, id)
}
