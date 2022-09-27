package quiz

import (
	"encoding/json"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

const (
	Buzzed = iota
)

const (
	StartNewQuestion = iota
	SetWrongAnswer
	SetCorrectAnswer
	RemoveUser
)

const (
	Init = iota
	UserEntered
	UserLeft
	UserStatus
	NewRound
)

const (
	none = iota
	active
	blocked
)

type userRequest struct {
	Action     int8   `json:"action"`
	Payload    string `json:"payload"`
	Connection *websocket.Conn
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
	Id      string `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"-"`
	Points  int32  `json:"points"`
	State   int8   `json:"state"`
}

type Communication struct {
	logger       zerolog.Logger
	userPool     *ConnectionPool
	adminPool    *ConnectionPool
	overviewPool *ConnectionPool
	userState    map[*websocket.Conn]*user
	isBuzzed     *websocket.Conn
}

func NewCommunication(logger zerolog.Logger) *Communication {
	userReceived := make(chan recMsg)
	userConnectionPool := NewConnectionPool(logger, userReceived)
	adminReceived := make(chan recMsg)
	adminConnectionPool := NewConnectionPool(logger, adminReceived)
	overviewReceived := make(chan recMsg)
	overviewConnectionPool := NewConnectionPool(logger, overviewReceived)
	userState := make(map[*websocket.Conn]*user)

	communication := &Communication{
		logger:       logger,
		userPool:     userConnectionPool,
		adminPool:    adminConnectionPool,
		overviewPool: overviewConnectionPool,
		isBuzzed:     nil,
		userState:    userState,
	}

	go communication.run(userReceived, adminReceived)

	return communication
}

func (C *Communication) run(userReceived, adminReceived chan recMsg) {
	// todo maybe it is better to listen to each received channel on its own to handle admin messages instant
	for {
		select {
		case message := <-userReceived:
			C.logger.Debug().Str("user: ", string(message.Payload))
			userRequest, err := unmarshal(message.Payload, userRequest{})
			if err != nil {
				C.logger.Debug().Err(err)
				return
			}
			C.handleUserRequest(userRequest, message.Connection)
		case message := <-adminReceived:
			C.logger.Debug().Str("admin: ", string(message.Payload))
			adminRequest, err := unmarshal(message.Payload, adminRequest{})
			if err != nil {
				C.logger.Debug().Err(err)
				return
			}
			C.handleAdminRequest(adminRequest)
		}
	}
}

func (C *Communication) handleUserRequest(message userRequest, connection *websocket.Conn) {
	switch message.Action {
	case Buzzed:
		if C.isBuzzed != nil {
			return
		}
		user := C.userState[connection]
		C.isBuzzed = connection
		user.State = active
		C.sendUserStatus(user)
	}
}

func (C *Communication) handleAdminRequest(message adminRequest) {
	switch message.Action {
	case StartNewQuestion:
		C.isBuzzed = nil
		for _, user := range C.userState {
			if user.State != none {
				user.State = none
				C.sendUserStatus(user)
			}
		}
		C.sendNewRound(message.Payload)
	case SetWrongAnswer:
		for connection, user := range C.userState {
			if connection == C.isBuzzed {
				user.State = blocked
				user.Points -= 2
			} else {
				user.Points += 1
			}
			C.sendUserStatus(user)
		}
		C.isBuzzed = nil
	case SetCorrectAnswer:
		// add points to is buzzed user
		user := C.userState[C.isBuzzed]
		user.Points += 5
		C.sendUserStatus(user)
	case RemoveUser:
		// todo
	}
}

func (C *Communication) sendInitialUserData(connection *websocket.Conn, pool *ConnectionPool) {
	var users []*user
	for _, user := range C.userState {
		users = append(users, user)
	}
	response := response{
		Action:  Init,
		Payload: users,
	}
	pool.Broadcast(connection, response)
}

func (C *Communication) sendUserEntered(user *user) {
	response := response{
		Action:  UserEntered,
		Payload: user,
	}
	C.adminPool.BroadcastAll(response)
	C.userPool.BroadcastAll(response)
	C.overviewPool.BroadcastAll(response)
}

func (C *Communication) sendUserLeft(id string) {
	response := response{
		Action:  UserLeft,
		Payload: id,
	}
	C.adminPool.BroadcastAll(response)
	C.userPool.BroadcastAll(response)
	C.overviewPool.BroadcastAll(response)
}

func (C *Communication) sendUserStatus(user *user) {
	response := response{
		Action:  UserStatus,
		Payload: user,
	}
	C.adminPool.BroadcastAll(response)
	C.userPool.BroadcastAll(response)
	C.overviewPool.BroadcastAll(response)
}

func (C *Communication) sendNewRound(question string) {
	response := response{
		Action:  NewRound,
		Payload: question,
	}
	C.userPool.BroadcastAll(response)
	C.overviewPool.BroadcastAll(response)
}

func (C *Communication) RegisterUser(connection *websocket.Conn) {
	user, err := C.getUserFromConnection(connection)
	if err != nil {
		return
	}
	C.userState[connection] = user
	defer func() {
		C.sendUserLeft(user.Id)
		delete(C.userState, connection)
		C.userPool.Unregister(connection)
		connection.Close()
	}()
	// todo check for race condition
	C.sendUserEntered(user)
	C.sendInitialUserData(connection, C.userPool)
	C.userPool.Register(connection)
}

func (C *Communication) RegisterAdmin(connection *websocket.Conn) {
	defer func() {
		// todo disconnect all users || maybe later ;)
		C.adminPool.Unregister(connection)
		connection.Close()
	}()
	C.sendInitialUserData(connection, C.adminPool)
	C.adminPool.Register(connection)
}

func (C *Communication) RegisterOverview(connection *websocket.Conn) {
	defer func() {
		C.overviewPool.Unregister(connection)
		connection.Close()
	}()
	C.sendInitialUserData(connection, C.overviewPool)
	C.overviewPool.Register(connection)
}

func (C *Communication) getUserFromConnection(connection *websocket.Conn) (*user, error) {
	token := connection.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	user := &user{
		Id:      claims["id"].(string),
		Name:    claims["name"].(string),
		IsAdmin: claims["isAdmin"].(bool),
		Points:  0,
		State:   none,
	}

	return user, nil
}

func unmarshal[T any](message []byte, genericInterface T) (T, error) {
	err := json.Unmarshal(message, &genericInterface)
	if err != nil {
		return genericInterface, err
	}
	return genericInterface, nil
}
