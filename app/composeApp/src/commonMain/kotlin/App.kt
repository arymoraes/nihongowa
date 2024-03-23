import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Text
import androidx.compose.material.Card
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import org.jetbrains.compose.resources.ExperimentalResourceApi
import org.jetbrains.compose.ui.tooling.preview.Preview

import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.serialization.kotlinx.json.*

import kotlinx.serialization.json.*

import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.ui.unit.dp
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.sp
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items

import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material.Button
import androidx.compose.material.TextField
import androidx.compose.runtime.MutableState
import androidx.compose.runtime.remember

import model.Message

@OptIn(ExperimentalResourceApi::class)
@Composable
@Preview
fun App() {
    val messages = remember { mutableStateOf<List<Message>>(emptyList()) }
    val inputText = remember { mutableStateOf("") }

    MaterialTheme {
        LaunchedEffect(Unit) {
            messages.value = fetchData()
        }
        Column(Modifier.fillMaxSize()) {
            LazyColumn(modifier = Modifier.weight(1f)) {
                items(messages.value) { message ->
                    MessageCard(message = message)
                }
            }
            // Input field with a send button
            MessageInputField(inputText = inputText)
        }
    }
}


@Composable
fun MessageCard(message: Message) {
    Box(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 8.dp, vertical = 4.dp)
    ) {
        Card(
            modifier = Modifier.align(if (message.translation.isNotEmpty()) Alignment.CenterStart else Alignment.CenterEnd),
            shape = RoundedCornerShape(8.dp),
            elevation = 4.dp,
            backgroundColor = if (message.translation.isNotEmpty()) Color(0xFFE7FFDB) else Color(0xFFFFFFFF) // Greenish for bot, whiteish for user
        ) {
            Column(
                modifier = Modifier
                    .padding(16.dp)
                    .fillMaxWidth(),
                horizontalAlignment = if (message.translation.isNotEmpty()) Alignment.Start else Alignment.End
            ) {
                Text(
                    text = message.content,
                    fontSize = 16.sp,
                    fontWeight = FontWeight.Bold
                )
                if (message.translation.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(
                        text = "Translation: ${message.translation}",
                        fontSize = 14.sp,
                        fontWeight = FontWeight.Normal,
                        color = Color.Gray
                    )
                }
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = "Time: now",
                    fontSize = 12.sp,
                    color = Color.Gray
                )
            }
        }
    }
}

@Composable
fun MessageInputField(inputText: MutableState<String>) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(8.dp),
        verticalAlignment = Alignment.CenterVertically
    ) {
        TextField(
            value = inputText.value,
            onValueChange = { inputText.value = it },
            modifier = Modifier.weight(1f),
            placeholder = { Text("Type a message") },
            maxLines = 4
        )
        Button(
            onClick = {
                // Handle send message action here
                println("Send message: ${inputText.value}")
                // Clear the input field after sending
                inputText.value = ""
            },
            modifier = Modifier.padding(start = 8.dp)
        ) {
            Text("Send")
        }
    }
}

suspend fun fetchData(): List<Message> {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    try {
        val response: HttpResponse = client.get("http://192.168.1.71:1323/messages/582633d3-e87c-4ef1-9ff6-75f6c3c80751")

        return Json.decodeFromString<List<Message>>(response.bodyAsText())

    } catch (e: Exception) {
        println("Error: ${e.message}")
        return emptyList()
    }
}