import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material.MaterialTheme
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import io.ktor.client.*
import io.ktor.client.engine.cio.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.client.request.forms.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.*
import model.Message
import org.jetbrains.compose.resources.ExperimentalResourceApi
import org.jetbrains.compose.ui.tooling.preview.Preview

expect fun currentTimeMillis(): Long

@OptIn(ExperimentalResourceApi::class)
@Composable
@Preview
fun App() {
    val messages = remember { mutableStateOf<List<Message>>(emptyList()) }
    val inputText = remember { mutableStateOf("") }
    val lastRefresh = remember { mutableStateOf(currentTimeMillis()) }
    val listState = rememberLazyListState() // Step 1: Create LazyListState

    MaterialTheme {
        // Fetch messages
        LaunchedEffect(lastRefresh.value) {
            messages.value = fetchData()
        }

        // Scroll to the bottom when the list of messages changes
        LaunchedEffect(messages.value.size) {
            if (messages.value.isNotEmpty()) {
                listState.scrollToItem(messages.value.size - 1)
            }
        }

        Column(Modifier.fillMaxSize()) {
            LazyColumn(state = listState, modifier = Modifier.weight(1f).fillMaxWidth()) {
                items(messages.value) { message ->
                    MessageCard(message = message)
                }
            }
            MessageInputField(inputText = inputText, onMessageSent = {
                lastRefresh.value = currentTimeMillis()
            })
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

suspend fun sendMessage(message: String): Boolean {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    return try {
        val response: HttpResponse = client.submitForm(
            url = "http://192.168.1.71:1323/messages/582633d3-e87c-4ef1-9ff6-75f6c3c80751",
            formParameters = Parameters.build {
                append("message", message)
            }
        )
        response.status == HttpStatusCode.OK
    } catch (e: Exception) {
        println("Error sending message: ${e.message}")
        false
    } finally {
        client.close()
    }
}
