package quiz

import (
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

type ConnectionPool struct {
	logger     zerolog.Logger
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan string
}

func NewConnectionPool(logger zerolog.Logger) *ConnectionPool {
	connectionPool := &ConnectionPool{logger: logger}
	connectionPool.clients = make(map[*websocket.Conn]bool) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	connectionPool.register = make(chan *websocket.Conn)
	connectionPool.unregister = make(chan *websocket.Conn)
	connectionPool.broadcast = make(chan string)

	return connectionPool
}

func (C *ConnectionPool) Run() {
	for {
		select {
		case connection := <-C.register:
			C.clients[connection] = true
			C.logger.Debug().Msg("connection registered")

		case connection := <-C.unregister:
			// Remove the client from the hub
			delete(C.clients, connection)

			C.logger.Debug().Msg("connection unregistered")

		case message := <-C.broadcast:
			C.logger.Debug().Str("send message", message)

			// Send the message to all clients
			for connection := range C.clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					C.logger.Error().Str("write error", err.Error())

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(C.clients, connection)
				}
			}
		}
	}
}

func (C *ConnectionPool) Register(connection *websocket.Conn) {
	C.register <- connection
}

func (C *ConnectionPool) Unregister(connection *websocket.Conn) {
	C.unregister <- connection
}

func (C *ConnectionPool) Broadcast(message string) {
	C.broadcast <- message
}
