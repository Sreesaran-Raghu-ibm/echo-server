package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for testing
	},
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("Client connected: %s", r.RemoteAddr)

	// Send welcome message
	// welcomeMsg := fmt.Sprintf("Welcome! Connected to echo server at %s", r.Host)
	// if err := conn.WriteMessage(websocket.TextMessage, []byte(welcomeMsg)); err != nil {
	// 	log.Printf("Write error: %v", err)
	// 	return
	// }

	// Echo loop - read and echo back messages
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error: %v", err)
			}
			break
		}

		// Log received message
		if messageType == websocket.TextMessage {
			log.Printf("Received (text): %s", message)
		} else {
			log.Printf("Received (binary): %d bytes", len(message))
		}

		// Echo the message back
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Printf("Write error: %v", err)
			break
		}
	}

	log.Printf("Client disconnected: %s", r.RemoteAddr)
}

// Serve HTML test page
func handleHome(w http.ResponseWriter, r *http.Request) {
	html := `<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Echo Test</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 50px auto; padding: 20px; }
        #messages { border: 1px solid #ccc; height: 300px; overflow-y: scroll; padding: 10px; margin: 20px 0; background: #f5f5f5; }
        .message { margin: 5px 0; padding: 5px; }
        .sent { color: blue; }
        .received { color: green; }
        .error { color: red; }
        .info { color: gray; font-style: italic; }
        input[type="text"] { width: 70%; padding: 10px; }
        button { padding: 10px 20px; margin-left: 10px; cursor: pointer; }
        .status { padding: 10px; margin: 10px 0; border-radius: 5px; }
        .connected { background: #d4edda; color: #155724; }
        .disconnected { background: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <h1>WebSocket Echo Server Test</h1>
    
    <div id="status" class="status disconnected">Disconnected</div>
    
    <button id="connectBtn" onclick="connect()">Connect</button>
    <button id="disconnectBtn" onclick="disconnect()" disabled>Disconnect</button>
    
    <div id="messages"></div>
    
    <div>
        <input type="text" id="messageInput" placeholder="Type a message..." onkeypress="handleKeyPress(event)" disabled>
        <button onclick="sendMessage()" id="sendBtn" disabled>Send</button>
    </div>

    <script>
        let ws = null;
        const messagesDiv = document.getElementById('messages');
        const messageInput = document.getElementById('messageInput');
        const statusDiv = document.getElementById('status');
        const connectBtn = document.getElementById('connectBtn');
        const disconnectBtn = document.getElementById('disconnectBtn');
        const sendBtn = document.getElementById('sendBtn');

        function addMessage(text, className) {
            const msg = document.createElement('div');
            msg.className = 'message ' + className;
            msg.textContent = text;
            messagesDiv.appendChild(msg);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }

        function updateStatus(connected) {
            if (connected) {
                statusDiv.textContent = 'Connected';
                statusDiv.className = 'status connected';
                connectBtn.disabled = true;
                disconnectBtn.disabled = false;
                messageInput.disabled = false;
                sendBtn.disabled = false;
            } else {
                statusDiv.textContent = 'Disconnected';
                statusDiv.className = 'status disconnected';
                connectBtn.disabled = false;
                disconnectBtn.disabled = true;
                messageInput.disabled = true;
                sendBtn.disabled = true;
            }
        }

        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';
            
            addMessage('Connecting to ' + wsUrl + '...', 'info');
            
            ws = new WebSocket(wsUrl);

            ws.onopen = function(event) {
                addMessage('Connected to server!', 'info');
                updateStatus(true);
            };

            ws.onmessage = function(event) {
                addMessage('← ' + event.data, 'received');
            };

            ws.onerror = function(error) {
                addMessage('Error: ' + error, 'error');
            };

            ws.onclose = function(event) {
                addMessage('Disconnected from server', 'info');
                updateStatus(false);
                ws = null;
            };
        }

        function disconnect() {
            if (ws) {
                ws.close();
            }
        }

        function sendMessage() {
            const message = messageInput.value.trim();
            if (message && ws && ws.readyState === WebSocket.OPEN) {
                ws.send(message);
                addMessage('→ ' + message, 'sent');
                messageInput.value = '';
            }
        }

        function handleKeyPress(event) {
            if (event.key === 'Enter') {
                sendMessage();
            }
        }

        // Auto-connect on page load
        window.onload = function() {
            connect();
        };
    </script>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	// Routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/ws", handleWebSocket)

	port := "8080"
	log.Printf("WebSocket Echo Server starting on port %s", port)
	log.Printf("Open http://localhost:%s in your browser", port)
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}

// Made with Bob
