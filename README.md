# WebSocket Echo Server Example

A simple, standalone WebSocket echo server with a built-in web interface for testing.

## Features

- WebSocket server that echoes back any message it receives
- Built-in HTML test page with interactive UI
- Supports both text and binary messages
- Connection status indicator
- Message history display
- Auto-connect on page load

## Running the Server

```bash
cd examples
go run websocket-echo-server.go
```

The server will start on port 8080.

## Testing the Server

### Option 1: Use the Built-in Web Interface

1. Open your browser and navigate to: `http://localhost:8080`
2. The page will automatically connect to the WebSocket server
3. Type messages in the input field and click "Send" (or press Enter)
4. Watch your messages being echoed back in real-time

### Option 2: Use Command Line Tools

**Using `websocat` (install from https://github.com/vi/websocat):**
```bash
websocat ws://localhost:8080/ws
```

**Using `wscat` (install via npm: `npm install -g wscat`):**
```bash
wscat -c ws://localhost:8080/ws
```

### Option 3: Use JavaScript in Browser Console

```javascript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onopen = () => {
    console.log('Connected!');
    ws.send('Hello, Server!');
};

ws.onmessage = (event) => {
    console.log('Received:', event.data);
};

ws.onerror = (error) => {
    console.error('Error:', error);
};

ws.onclose = () => {
    console.log('Disconnected');
};
```

## How It Works

1. **HTTP Upgrade**: The server upgrades HTTP connections to WebSocket protocol
2. **Welcome Message**: Upon connection, sends a welcome message to the client
3. **Echo Loop**: Continuously reads messages and echoes them back
4. **Logging**: Logs all connections, messages, and disconnections to console

## Code Structure

- **`upgrader`**: Configures WebSocket upgrade with CORS enabled
- **`handleWebSocket()`**: Handles WebSocket connections and message echoing
- **`handleHome()`**: Serves the HTML test interface
- **`main()`**: Sets up routes and starts the HTTP server

## Example Output

Server console:
```
2025/11/18 15:05:23 WebSocket Echo Server starting on port 8080
2025/11/18 15:05:23 Open http://localhost:8080 in your browser
2025/11/18 15:05:30 Client connected: [::1]:54321
2025/11/18 15:05:35 Received (text): Hello, Server!
2025/11/18 15:05:40 Received (text): Testing echo...
2025/11/18 15:05:45 Client disconnected: [::1]:54321
```

## Dependencies

This example uses the Gorilla WebSocket library:
```bash
go get github.com/gorilla/websocket
```

## Customization

You can modify the server by:
- Changing the port in `main()` function
- Adding message processing logic in the echo loop
- Implementing custom message types or protocols
- Adding authentication or rate limiting
- Storing message history