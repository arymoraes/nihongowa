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

	messages, err := services.GetMessagesFromConversation(id)

	// If it is an empty array its ok
	if messages == nil {
		return c.JSON(http.StatusNotFound, ErrorResponse{Message: "Conversation not found"})
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, messages)
}

// e.GET("/conversations", getLastConversations)
func GetLastConversations(c echo.Context) error {
	conversations, err := services.GetLastConversations()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, conversations)
}

// e.POST("/messages/:conversation_id", postMessageToConversation)
func PostMessageToConversation(c echo.Context) error {
	id := c.Param("conversation_id")

	content := c.FormValue("message")

	message := models.Message{
		Content: content,
	}

	reply, err := services.PostMessageToConversation(id, message)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, reply)
}

// e.POST("/conversations", postConversation)
func PostConversation(c echo.Context) error {
	type PostConversationResponse struct {
		ConversationID string `json:"conversation_id"`
		AssistantName  string `json:"assistant_name"`
	}

	conversation, err := services.CreateConversationScenario()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, PostConversationResponse{ConversationID: conversation.ID.String(), AssistantName: conversation.AssistantName})
}
