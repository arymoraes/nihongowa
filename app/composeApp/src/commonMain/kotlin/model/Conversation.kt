package model

import kotlinx.serialization.Serializable

@Serializable

data class Conversation(
    val id: String,
    val messages: List<Message>,
    val createdAt: String,
    val updatedAt: String,
)