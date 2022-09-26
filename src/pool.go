package quiz

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"strconv"
)

type ConnectionPool struct {
	logger     zerolog.Logger
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan interface{}
	received   chan []byte
}

func NewConnectionPool(logger zerolog.Logger, received chan []byte) *ConnectionPool {
	connectionPool := &ConnectionPool{logger: logger}
	connectionPool.clients = make(map[*websocket.Conn]bool) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	connectionPool.register = make(chan *websocket.Conn)
	connectionPool.unregister = make(chan *websocket.Conn)
	connectionPool.broadcast = make(chan interface{})
	connectionPool.received = received

	go connectionPool.run()

	return connectionPool
}

func (C *ConnectionPool) run() {
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
			C.logger.Debug().Msgf("send message to all: ", message)
			// Send the message to all clients
			for connection := range C.clients {
				C.sendMessage(connection, message)
			}
		}
	}
}

func (C *ConnectionPool) Register(connection *websocket.Conn) {
	// Register the client
	C.register <- connection

	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				C.logger.Error().Str("read error", err.Error())
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			C.received <- message
		} else {
			C.logger.Debug().Str("websocket message received of type", strconv.Itoa(messageType))
		}
	}
}

func (C *ConnectionPool) Unregister(connection *websocket.Conn) {
	C.unregister <- connection
	connection.Close()
}

func (C *ConnectionPool) Broadcast(connection *websocket.Conn, message string) {
	C.logger.Debug().Str("send message to client: ", message)
	C.sendMessage(connection, message)
}

func (C *ConnectionPool) BroadcastAll(message interface{}) {
	C.broadcast <- message
}

func (C *ConnectionPool) sendMessage(connection *websocket.Conn, message interface{}) {
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}
	if err = connection.WriteMessage(websocket.TextMessage, payload); err != nil {
		C.logger.Error().Str("write error", err.Error())

		connection.WriteMessage(websocket.CloseMessage, []byte{})
		connection.Close()
		delete(C.clients, connection)
	}
}
