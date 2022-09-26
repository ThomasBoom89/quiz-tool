package quiz

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
)

const (
	Buzzered = iota
)

const (
	StartNewQuestion = iota
	SetWrongAnswer
	SetCorrectAnswer
	RemoveUser
)

const (
	UserEntered = iota
	UserLeft
	UserStatus
	NewRound
)

type userRequest struct {
	Action  int8   `json:"action"`
	Payload string `json:"payload"`
}

type adminRequest struct {
	Action  int8   `json:"action"`
	Payload string `json:"payload"`
}

type response struct {
	Action  int8        `json:"action"`
	Payload interface{} `json:"payload"`
}

type user struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Points int32  `json:"points"`
	State  int8   `json:"state"`
}

type Communication struct {
	logger       zerolog.Logger
	userPool     *ConnectionPool
	adminPool    *ConnectionPool
	overviewPool *ConnectionPool
}

func NewCommunication(logger zerolog.Logger) *Communication {
	userReceived := make(chan []byte)
	userConnectionPool := NewConnectionPool(logger, userReceived)
	adminReceived := make(chan []byte)
	adminConnectionPool := NewConnectionPool(logger, adminReceived)
	overviewReceived := make(chan []byte)
	overviewConnectionPool := NewConnectionPool(logger, overviewReceived)

	communication := &Communication{
		logger:       logger,
		userPool:     userConnectionPool,
		adminPool:    adminConnectionPool,
		overviewPool: overviewConnectionPool,
	}

	go communication.run(userReceived, adminReceived)

	return communication
}

func (C *Communication) run(userReceived, adminReceived chan []byte) {
	// todo maybe it is better to listen to each received channel on its own to handle admin messages instant
	for {
		select {
		case message := <-userReceived:
			C.logger.Debug().Str("user: ", string(message))
			userRequest, err := unmarshal(message, userRequest{})
			if err != nil {
				C.logger.Debug().Err(err)
				return
			}
			C.handleUserRequest(userRequest)
		case message := <-adminReceived:
			C.logger.Debug().Str("admin: ", string(message))
			adminRequest, err := unmarshal(message, adminRequest{})
			if err != nil {
				C.logger.Debug().Err(err)
				return
			}
			C.handleAdminRequest(adminRequest)
		}
	}
}

func (C *Communication) handleUserRequest(message userRequest) {
	switch message.Action {
	case Buzzered:
		// todo buzzer action dinge
	}
}

func (C *Communication) handleAdminRequest(message adminRequest) {
	switch message.Action {
	case StartNewQuestion:
	case SetWrongAnswer:
	case SetCorrectAnswer:
	case RemoveUser:

	}
}

func (C *Communication) sendNewRoundResponse(question string) {
	response := response{
		Action:  NewRound,
		Payload: question,
	}
	C.adminPool.BroadcastAll(response)
}

func (C *Communication) sendUserResponse(connection *websocket.Conn, action int8) {
	// todo get user data by connection
	user := user{
		Id:     "",
		Name:   "",
		Points: 0,
		State:  0,
	}
	response := response{
		Action:  action,
		Payload: user,
	}
	C.adminPool.BroadcastAll(response)
}

func (C *Communication) RegisterUser(connection *websocket.Conn, jwt string) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		// todo send message to all other connection that one player is gone
		C.userPool.Unregister(connection)
		connection.Close()
	}()
	// todo send initial message to connection like other users etc.
	C.userPool.Register(connection)
}

func (C *Communication) RegisterAdmin(connection *websocket.Conn, jwt string) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		// todo disconnect all users
		C.adminPool.Unregister(connection)
		connection.Close()
	}()
	// todo send initial message to connection like other users etc.
	C.adminPool.Register(connection)
}

func (C *Communication) RegisterOverview(connection *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	defer func() {
		C.overviewPool.Unregister(connection)
		connection.Close()
	}()
	// todo send initial message to connection like other users etc.
	C.overviewPool.Register(connection)
}

func unmarshal[T any](message []byte, genericInterface T) (T, error) {
	err := json.Unmarshal(message, &genericInterface)
	if err != nil {
		return genericInterface, err
	}
	return genericInterface, nil
}
