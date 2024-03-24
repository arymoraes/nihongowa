// import androidx.compose.foundation.layout.Column
// import androidx.compose.foundation.layout.fillMaxSize
// import androidx.compose.foundation.layout.fillMaxWidth
// import androidx.compose.foundation.lazy.LazyColumn
// import androidx.compose.foundation.lazy.items
// import androidx.compose.foundation.lazy.rememberLazyListState
// import androidx.compose.material.MaterialTheme
// import androidx.compose.runtime.Composable
// import androidx.compose.runtime.LaunchedEffect
// import androidx.compose.runtime.mutableStateOf
// import androidx.compose.runtime.remember
// import androidx.compose.ui.Modifier
// import model.Message

// object ConversationScreen : Screen {

//     @Composable
//     override fun Content() {
//         ConversationScreen("http://192.168.1.71:1323/messages/cf45fcce-a794-47b6-aae4-cf26c0bce61a")
//     }
// }

// @Composable
// fun ConversationScreen(conversationId: String) {
//     val messages = remember { mutableStateOf<List<Message>>(emptyList()) }
//     val inputText = remember { mutableStateOf("") }
//     val lastRefresh = remember { mutableStateOf(currentTimeMillis()) }
//     val listState = rememberLazyListState() // Step 1: Create LazyListState

//     MaterialTheme {
//         // Fetch messages
//         LaunchedEffect(lastRefresh.value) {
//             messages.value = fetchData()
//         }

//         // Scroll to the bottom when the list of messages changes
//         LaunchedEffect(messages.value.size) {
//             if (messages.value.isNotEmpty()) {
//                 listState.scrollToItem(messages.value.size - 1)
//             }
//         }

//         Column(Modifier.fillMaxSize()) {
//             LazyColumn(state = listState, modifier = Modifier.weight(1f).fillMaxWidth()) {
//                 items(messages.value) { message ->
//                     MessageCard(message = message)
//                 }
//             }
//             MessageInputField(inputText = inputText, onMessageSent = {
//                 lastRefresh.value = currentTimeMillis()
//             })
//         }
//     }
// }
