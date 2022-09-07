package quiz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"strconv"
)

type Overview struct {
	router     fiber.Router
	logger     zerolog.Logger
	clients    map[*websocket.Conn]bool
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	broadcast  chan string
}

func NewOverview(router fiber.Router, logger zerolog.Logger) *Overview {
	overview := &Overview{router: router, logger: logger}
	overview.clients = make(map[*websocket.Conn]bool) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	overview.register = make(chan *websocket.Conn)
	overview.broadcast = make(chan string)
	overview.unregister = make(chan *websocket.Conn)
	overview.attachRoutes()

	return overview
}

func (U *Overview) Run() {
	for {
		select {
		case connection := <-U.register:
			U.clients[connection] = true
			U.logger.Debug().Msg("connection registered")

		case message := <-U.broadcast:
			U.logger.Debug().Str("send message", message)

			// Send the message to all clients
			for connection := range U.clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					U.logger.Error().Str("write error", err.Error())

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(U.clients, connection)
				}
			}

		case connection := <-U.unregister:
			// Remove the client from the hub
			delete(U.clients, connection)

			U.logger.Debug().Msg("connection unregistered")
		}
	}
}

func (U *Overview) attachRoutes() {
	U.router.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	U.router.Get("/ws", websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			U.unregister <- c
			c.Close()
		}()

		// Register the client
		U.register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					U.logger.Error().Str("read error", err.Error())
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				U.broadcast <- string(message)
			} else {
				U.logger.Error().Str("websocket message received of type", strconv.Itoa(messageType))
			}
		}
	}))
}
