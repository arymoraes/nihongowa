package screen.conversation

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material.Icon
import androidx.compose.material.IconButton
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Text
import androidx.compose.material.TopAppBar
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import cafe.adriel.voyager.core.screen.Screen
import cafe.adriel.voyager.navigator.LocalNavigator
import currentTimeMillis
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.request.get
import io.ktor.client.statement.HttpResponse
import io.ktor.client.statement.bodyAsText
import io.ktor.serialization.kotlinx.json.json
import kotlinx.serialization.json.Json
import model.Message
import org.jetbrains.compose.resources.ExperimentalResourceApi


data class ConversationScreen(val conversationId: String, val assistantName: String): Screen {
    @OptIn(ExperimentalResourceApi::class)
    @Composable
    override fun Content() {
        val navigator = LocalNavigator.current // Get the current navigator instance
        val messages = remember { mutableStateOf<List<Message>>(emptyList()) }
        val inputText = remember { mutableStateOf("") }
        val lastRefresh = remember { mutableStateOf(currentTimeMillis()) }
        val listState = rememberLazyListState()

        MaterialTheme {
            Column(Modifier.fillMaxSize()) {
                // TopAppBar as a header
                TopAppBar(
                    title = { Text(text = "$assistantName") },
                    backgroundColor = Color.Blue, // Customize as needed
                    contentColor = Color.White,
                    navigationIcon = {
                        IconButton(onClick = { navigator?.pop() }) {
                            // Use Material Icons directly
                            Icon(Icons.Filled.ArrowBack, contentDescription = "Back")
                        }
                    }
                )

                // Fetch messages
                LaunchedEffect(lastRefresh.value) {
                    messages.value = fetchData(conversationId)
                }

                // Scroll to the bottom when the list of messages changes
                LaunchedEffect(messages.value.size) {
                    if (messages.value.isNotEmpty()) {
                        listState.scrollToItem(messages.value.size - 1)
                    }
                }

                LazyColumn(state = listState, modifier = Modifier.weight(1f).fillMaxWidth()) {
                    items(messages.value) { message ->
                        MessageCard(message = message)
                    }
                }

                MessageInputField(inputText = inputText, onMessageSent = {
                    lastRefresh.value = currentTimeMillis()
                }, conversationId = conversationId)
            }
        }
    }
}

suspend fun fetchData(conversationId: String): List<Message> {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    try {
        val response: HttpResponse = client.get("http://192.168.1.71:1323/messages/$conversationId")

        return Json.decodeFromString<List<Message>>(response.bodyAsText())

    } catch (e: Exception) {
        println("Error: ${e.message}")
        return emptyList()
    }
}