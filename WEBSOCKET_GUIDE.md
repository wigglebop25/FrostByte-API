# WebSocket Integration Guide for Android (Kotlin)

This guide explains how to connect your Android app to the FrostByte API WebSocket for real-time order updates.

## 1. Connection Details

* **Endpoint:** `/ws`
* **Full URL:** `wss://<your-server-ip-or-domain>/ws`
 * *Production:* `wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws`
 * *Localhost (Emulator):* `ws://10.0.2.2:8080/ws`
* **Method:** `GET`
* **Authentication:** **Query Parameter**
 * You must pass the JWT token in the URL query string: `?token=<your_jwt_token>`

## 2. Dependencies (Recommended)

Use **OkHttp** for robust WebSocket handling in Android. Add this to your `build.gradle.kts` (Module level):

```kotlin
implementation("com.squareup.okhttp3:okhttp:4.12.0")
```

## 3. Implementation Example

### Step A: Define the WebSocket Listener

Create a class that extends `WebSocketListener` to handle incoming events.

```kotlin
import okhttp3.Response
import okhttp3.WebSocket
import okhttp3.WebSocketListener
import android.util.Log

class OrderWebSocketListener(private val onMessageReceived: (String) -> Unit) : WebSocketListener() {

 override fun onOpen(webSocket: WebSocket, response: Response) {
 super.onOpen(webSocket, response)
 Log.d("FrostByteWS", "Connection Established: ${response.message}")
 }

 override fun onMessage(webSocket: WebSocket, text: String) {
 super.onMessage(webSocket, text)
 Log.d("FrostByteWS", "Message Received: $text")
 // Pass the message to the UI or Data layer
 onMessageReceived(text)
 }

 override fun onClosing(webSocket: WebSocket, code: Int, reason: String) {
 super.onClosing(webSocket, code, reason)
 webSocket.close(1000, null)
 Log.d("FrostByteWS", "Connection Closing: $code / $reason")
 }

 override fun onFailure(webSocket: WebSocket, t: Throwable, response: Response?) {
 super.onFailure(webSocket, t, response)
 Log.e("FrostByteWS", "Error: ${t.message}")
 }
}
```

### Step B: Connect to the Server

Initialize the connection, ensuring you pass the JWT token retrieved from your login flow.

```kotlin
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.WebSocket

class WebSocketManager {
 private var webSocket: WebSocket? = null
 private val client = OkHttpClient()

 fun connect(jwtToken: String) {
 // 1. Construct the URL with the Token
 val wsUrl = "wss://frostbyte-api.southeastasia.cloudapp.azure.com/ws?token=$jwtToken"

 // 2. Create the Request
 val request = Request.Builder()
 .url(wsUrl)
 .build()

 // 3. Create the Listener
 val listener = OrderWebSocketListener { message ->
 // Handle the JSON message here (e.g., parse it with Gson/Moshi)
 println("New Order Update: $message")
 }

 // 4. Start the connection
 webSocket = client.newWebSocket(request, listener)
 }

 fun close() {
 webSocket?.close(1000, "App closed")
 }
}
```

## 4. Expected Message Format

The server sends JSON messages when order statuses change.

**Example Payload:**
```json
{
 "order_id": 123,
 "status": "READY", 
 "updated_at": "2024-01-13T10:00:00Z"
}
```

## 5. Best Practices

1. **Lifecycle Awareness:** Connect in `onResume()` (or a ViewModel) and close the connection in `onPause()` or `onCleared()` to save battery.
2. **Reconnection Logic:** WebSockets can drop. Implement a retry mechanism with exponential backoff in the `onFailure` method of your listener.
3. **Security:** Ensure you use `wss://` (Secure WebSocket) in production to encrypt the traffic.
