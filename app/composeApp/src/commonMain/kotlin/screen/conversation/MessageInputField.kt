package screen.conversation

import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.Button
import androidx.compose.material.Text
import androidx.compose.material.TextField
import androidx.compose.runtime.Composable
import androidx.compose.runtime.MutableState
import androidx.compose.runtime.rememberCoroutineScope
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import io.ktor.client.HttpClient
import io.ktor.client.engine.cio.CIO
import io.ktor.client.plugins.contentnegotiation.ContentNegotiation
import io.ktor.client.request.forms.submitForm
import io.ktor.client.statement.HttpResponse
import io.ktor.http.HttpStatusCode
import io.ktor.http.Parameters
import io.ktor.serialization.kotlinx.json.json
import kotlinx.coroutines.launch

@Composable
fun MessageInputField(inputText: MutableState<String>, onMessageSent: () -> Unit, conversationId: String) {
    val coroutineScope = rememberCoroutineScope()

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
                coroutineScope.launch {
                    val success = sendMessage(inputText.value, conversationId)
                    if (success) {
                        onMessageSent()
                        inputText.value = "" // Clear the input field after sending
                    }
                }
            },
            modifier = Modifier.padding(start = 8.dp)
        ) {
            Text("Send")
        }
    }
}

suspend fun sendMessage(message: String, conversationId: String): Boolean {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    return try {
        val response: HttpResponse = client.submitForm(
            url = "http://192.168.1.71:1323/messages/$conversationId",
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
