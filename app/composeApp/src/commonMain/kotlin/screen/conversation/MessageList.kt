package screen.conversation

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