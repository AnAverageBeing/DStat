package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	Addr     string
	server   *http.Server
	upgrader *websocket.Upgrader
	clients  map[*websocket.Conn]struct{}
	mu       sync.Mutex
}

func NewWebSocketServer(addr string) *WebSocketServer {
	return &WebSocketServer{
		Addr: addr,
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clients: make(map[*websocket.Conn]struct{}),
	}
}

func (ws *WebSocketServer) Start() {
	ws.server = &http.Server{
		Addr:    ws.Addr,
		Handler: http.HandlerFunc(ws.handleWebSocket),
	}
	go func() {
		if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("WebSocket server ListenAndServe: %v", err)
		}
	}()
}

func (ws *WebSocketServer) Stop() {
	if ws.server != nil {
		if err := ws.server.Close(); err != nil {
			log.Fatalf("WebSocket server Close: %v", err)
		}
	}
}

func (ws *WebSocketServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	ws.mu.Lock()
	ws.clients[conn] = struct{}{}
	ws.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	ws.mu.Lock()
	delete(ws.clients, conn)
	ws.mu.Unlock()
}

func (ws *WebSocketServer) Broadcast(data []byte) {
	ws.mu.Lock()
	defer ws.mu.Unlock()

	for conn := range ws.clients {
		err := conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			delete(ws.clients, conn)
		}
	}
}
