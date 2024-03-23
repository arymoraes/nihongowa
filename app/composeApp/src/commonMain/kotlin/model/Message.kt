package model

import kotlinx.serialization.*
import kotlinx.serialization.json.*

@Serializable 
data class Message(
    val content: String,
    val translation: String,
    val wordByWordTranslation: List<String>?,
    val createdAt: String,
    val updatedAt: String,
)