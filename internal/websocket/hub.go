package websocket

import (
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

// TokenValidator interface to avoid import cycle with service package
type TokenValidator interface {
	ValidateToken(tokenString string) (*jwt.Token, error)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for this implementation
	},
}

// BroadcastMessage defines the payload and targeting rules
type BroadcastMessage struct {
	Payload interface{}
	UserID  uint     // Target specific user (0 = broadcast to all if roles empty)
	Roles   []string // Target specific roles (empty = no role restriction)
}

type Hub struct {
	clients     map[*Client]bool
	broadcast   chan BroadcastMessage
	register    chan *Client
	unregister  chan *Client
	authService TokenValidator
}

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	UserID uint
	Roles  []string
}

func NewHub(authService TokenValidator) *Hub {
	return &Hub{
		broadcast:   make(chan BroadcastMessage),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		authService: authService,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Printf("Client connected: UserID=%d, Roles=%v", client.UserID, client.Roles)
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Printf("Client disconnected: UserID=%d", client.UserID)
			}
		case message := <-h.broadcast:
			payloadBytes, err := json.Marshal(message.Payload)
			if err != nil {
				log.Printf("Error marshaling broadcast message: %v", err)
				continue
			}

			for client := range h.clients {
				shouldSend := false

				// Logic:
				// 1. If Roles are specified, client must have at least one of the roles (OR logic for roles)
				// 2. If UserID is specified, client must match UserID
				// 3. If neither specified, broadcast to all (Global)

				// Check Roles (e.g., Admin, Cashier)
				if len(message.Roles) > 0 {
					for _, role := range message.Roles {
						for _, clientRole := range client.Roles {
							if role == clientRole {
								shouldSend = true
								break
							}
						}
						if shouldSend {
							break
						}
					}
				}

				// Check UserID (Specific Customer)
				// Note: If both Roles and UserID are present, it's an OR condition?
				// Usually: "Send to User X OR Admins". Let's assume inclusive OR.
				if !shouldSend && message.UserID > 0 {
					if client.UserID == message.UserID {
						shouldSend = true
					}
				}

				// If no specific targets, it's a global broadcast
				if len(message.Roles) == 0 && message.UserID == 0 {
					shouldSend = true
				}

				if shouldSend {
					select {
					case client.send <- payloadBytes:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}

// Broadcast sends a message to specific targets
func (h *Hub) Broadcast(payload interface{}, userID uint, roles []string) {
	h.broadcast <- BroadcastMessage{
		Payload: payload,
		UserID:  userID,
		Roles:   roles,
	}
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	// 1. Extract Token
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	// 2. Validate Token
	token, err := h.authService.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
		return
	}
	userID := uint(userIDFloat)

	var roles []string
	if rolesClaim, ok := claims["roles"].([]interface{}); ok {
		for _, r := range rolesClaim {
			if rStr, ok := r.(string); ok {
				roles = append(roles, rStr)
			}
		}
	}

	// 3. Upgrade Connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to websocket: %v", err)
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		UserID: userID,
		Roles:  roles,
	}
	h.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	
	// Configure connection limits
	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { 
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil 
	})

	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
