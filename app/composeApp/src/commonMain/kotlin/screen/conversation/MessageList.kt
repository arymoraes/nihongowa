package screen.conversation

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.Card
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import model.Message

@Composable
fun MessageList(messages: List<Message>) {
    LazyColumn(modifier = Modifier.fillMaxWidth()) {
        items(messages) { message ->
            MessageCard(message = message)
        }
    }
}

@Composable
fun MessageCard(message: Message) {
    // Track expanded state
    val isExpanded = remember { mutableStateOf(false) }

    Box(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 8.dp, vertical = 4.dp)
            .clickable { isExpanded.value = !isExpanded.value } // Toggle expanded state on click
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
                // Show translation and romanji only when expanded
                if (isExpanded.value && message.translation.isNotEmpty()) {
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(
                        text = "Translation: ${message.translation}",
                        fontSize = 14.sp,
                        fontWeight = FontWeight.Normal,
                        color = Color.Gray
                    )
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(
                        text = "Romanji: ${message.romanji}", // Assuming the property name is romanji
                        fontSize = 14.sp,
                        fontWeight = FontWeight.Normal,
                        color = Color.Gray
                    )
                }
                Spacer(modifier = Modifier.height(4.dp))
                Text(
                    text = "Time: now", // Consider updating this to show the actual time
                    fontSize = 12.sp,
                    color = Color.Gray
                )
            }
        }
    }
}