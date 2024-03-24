package screen.home

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.Button
import androidx.compose.material.Card
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import cafe.adriel.voyager.core.screen.Screen
import cafe.adriel.voyager.navigator.LocalNavigator
import currentTimeMillis
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.request.get
import io.ktor.client.request.post
import io.ktor.client.statement.HttpResponse
import io.ktor.client.statement.bodyAsText
import io.ktor.http.HttpStatusCode
import io.ktor.serialization.kotlinx.json.json
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import model.Conversation
import screen.conversation.ConversationScreen


class HomeScreen: Screen {
    @Composable
    override fun Content() {
        val navigator = LocalNavigator.current
        val conversations = remember { mutableStateOf<List<Conversation>>(emptyList()) }
        val lastRefresh = remember { mutableStateOf(currentTimeMillis()) }

        // Trigger re-fetching of data when lastRefresh changes
        LaunchedEffect(lastRefresh.value) {
            conversations.value = fetchData()
        }

        // Use Column to layout the button and the list vertically
        Column(modifier = Modifier.fillMaxSize()) {
            Button(
                onClick = {
                    // Use a coroutine to call createConversation and navigate
                    CoroutineScope(Dispatchers.Main).launch {
                        val newConversationId = createConversation()
                        newConversationId?.let {
                            navigator?.push(ConversationScreen(it))
                        } ?: run {
                            println("Failed to create a new conversation.")
                        }
                    }
                },
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(8.dp) // Add padding for aesthetics
            ) {
                Text("Create Conversation")
            }

            // Display conversations in a list
            LazyColumn(modifier = Modifier.weight(1f)) {
                items(conversations.value) { conversation ->
                    ConversationCard(conversation = conversation) {
                        navigator?.push(ConversationScreen(conversation.id))
                    }
                }
            }
        }
    }
}

@Composable
fun ConversationCard(conversation: Conversation, onClick: () -> Unit) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(8.dp)
            .clickable(onClick = onClick),
        elevation = 4.dp
    ) {
        Column(
            modifier = Modifier
                .padding(16.dp)
        ) {
            Text(text = "Conversation ID: ${conversation.id}")
            // You can add more details here as needed
        }
    }
}

suspend fun fetchData(): List<Conversation> {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    try {
        val response: HttpResponse = client.get("http://192.168.1.71:1323/conversations")

        return Json.decodeFromString<List<Conversation>>(response.bodyAsText())

    } catch (e: Exception) {
        println("Error: ${e.message}")
        return emptyList()
    }
}

@Serializable
data class PostConversationResponse(@SerialName("conversation_id") val ConversationID: String)

suspend fun createConversation(): String? {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    return try {
        val response: HttpResponse = client.post("http://192.168.1.71:1323/conversations")
        if (response.status == HttpStatusCode.OK) {
            val responseBody = Json.decodeFromString<PostConversationResponse>(response.bodyAsText())
            responseBody.ConversationID
        } else {
            println("Failed to create conversation, status code: ${response.status}")
            null
        }
    } catch (e: Exception) {
        println("Error creating conversation: ${e.message}")
        null
    } finally {
        client.close()
    }
}