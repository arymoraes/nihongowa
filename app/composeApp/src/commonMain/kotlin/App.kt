import androidx.compose.foundation.layout.*
import androidx.compose.material.MaterialTheme
import androidx.compose.runtime.*
import cafe.adriel.voyager.navigator.Navigator
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
import screen.home.HomeScreen

expect fun currentTimeMillis(): Long

@OptIn(ExperimentalResourceApi::class)
@Composable
@Preview
fun App() {
    MaterialTheme {
        Navigator(HomeScreen())
    }
}

suspend fun fetchData(): List<Message> {
    val client = HttpClient(CIO) {
        install(ContentNegotiation) {
            json()
        }
    }

    try {
        val response: HttpResponse = client.get("http://192.168.1.71:1323/messages/cf45fcce-a794-47b6-aae4-cf26c0bce61a")

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
            url = "http://192.168.1.71:1323/messages/cf45fcce-a794-47b6-aae4-cf26c0bce61a",
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
