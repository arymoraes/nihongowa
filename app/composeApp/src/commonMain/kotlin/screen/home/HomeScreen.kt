package screen.home

import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.lazy.rememberLazyListState
import androidx.compose.material.Button
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
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
import model.Conversation
import screen.conversation.ConversationScreen


class HomeScreen: Screen {
    @Composable
    override fun Content() {
      val navigator = LocalNavigator.current
      val conversations = remember { mutableStateOf<List<Conversation>>(emptyList()) }
        val lastRefresh = remember { mutableStateOf(currentTimeMillis()) }
        val listState = rememberLazyListState() // Step 1: Create LazyListState

        LaunchedEffect(lastRefresh.value) {
            conversations.value = fetchData()
            println("Conversations: ${conversations.value}")
            
        }
        Box(
            modifier = Modifier.fillMaxSize(),
            contentAlignment = Alignment.Center
        ) {
            Button(onClick = { navigator?.push(ConversationScreen("cf45fcce-a794-47b6-aae4-cf26c0bce61a")) }) {
                Text("Click me!")
            }
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