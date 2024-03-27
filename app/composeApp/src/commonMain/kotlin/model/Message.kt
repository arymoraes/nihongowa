package model

import kotlinx.serialization.Serializable

@Serializable 
data class Message(
    val id: String,
    val content: String,
    val translation: String,
    val romanji: String?,
    val conversation_id: String,
    val is_ai: Boolean,
    val user_message_translated: String?,
    val word_by_word_translation: List<String>?,
    val created_at: String,
    val updated_at: String,
)