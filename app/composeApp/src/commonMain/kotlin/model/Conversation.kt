package model

import kotlinx.serialization.Serializable

@Serializable

data class Conversation(
    val id: String,
    val messages: List<Message>?,
    val thread_id: String,
    val assistant_id: String,
    val scenario: String,
)