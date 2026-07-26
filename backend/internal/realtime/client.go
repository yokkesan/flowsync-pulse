package realtime

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	pingPeriod = pongWait * 9 / 10

	maxMessageSize = 1024

	sendBufferSize = 256
)

type Client struct {
	hub *Hub

	connection *websocket.Conn

	userID    uint64
	companyID uint64

	send chan []byte

	closeOnce sync.Once
}

func NewClient(
	hub *Hub,
	connection *websocket.Conn,
	userID uint64,
	companyID uint64,
) *Client {
	return &Client{
		hub:        hub,
		connection: connection,
		userID:     userID,
		companyID:  companyID,
		send: make(
			chan []byte,
			sendBufferSize,
		),
	}
}

func (c *Client) Start() {
	if c.hub == nil ||
		c.connection == nil ||
		c.userID == 0 ||
		c.companyID == 0 {
		c.closeConnection()
		return
	}

	c.hub.Register(c)

	go c.writePump()

	c.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.closeConnection()
	}()

	c.connection.SetReadLimit(
		maxMessageSize,
	)

	if err := c.connection.SetReadDeadline(
		time.Now().Add(pongWait),
	); err != nil {
		log.Printf(
			"failed to set realtime read deadline: user_id=%d company_id=%d error=%v",
			c.userID,
			c.companyID,
			err,
		)
		return
	}

	c.connection.SetPongHandler(
		func(string) error {
			return c.connection.SetReadDeadline(
				time.Now().Add(pongWait),
			)
		},
	)

	for {
		messageType, _, err :=
			c.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(
				err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
				websocket.CloseAbnormalClosure,
			) {
				log.Printf(
					"realtime client read failed: user_id=%d company_id=%d error=%v",
					c.userID,
					c.companyID,
					err,
				)
			}

			return
		}

		if messageType != websocket.TextMessage &&
			messageType != websocket.BinaryMessage {
			continue
		}

		// 現時点ではクライアントから受信する
		// WebSocketイベントは扱わない。
		// 接続維持とサーバーからの通知配信のみを担当する。
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(
		pingPeriod,
	)
	defer func() {
		ticker.Stop()
		c.closeConnection()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.connection.SetWriteDeadline(
				time.Now().Add(writeWait),
			); err != nil {
				log.Printf(
					"failed to set realtime write deadline: user_id=%d company_id=%d error=%v",
					c.userID,
					c.companyID,
					err,
				)
				return
			}

			if !ok {
				_ = c.connection.WriteMessage(
					websocket.CloseMessage,
					[]byte{},
				)
				return
			}

			if err := c.connection.WriteMessage(
				websocket.TextMessage,
				message,
			); err != nil {
				log.Printf(
					"failed to write realtime message: user_id=%d company_id=%d error=%v",
					c.userID,
					c.companyID,
					err,
				)
				return
			}

		case <-ticker.C:
			if err := c.connection.SetWriteDeadline(
				time.Now().Add(writeWait),
			); err != nil {
				log.Printf(
					"failed to set realtime ping deadline: user_id=%d company_id=%d error=%v",
					c.userID,
					c.companyID,
					err,
				)
				return
			}

			if err := c.connection.WriteMessage(
				websocket.PingMessage,
				nil,
			); err != nil {
				return
			}
		}
	}
}

func (c *Client) enqueue(
	message []byte,
) bool {
	select {
	case c.send <- message:
		return true

	default:
		log.Printf(
			"realtime client send buffer full: user_id=%d company_id=%d",
			c.userID,
			c.companyID,
		)
		return false
	}
}

func (c *Client) closeSend() {
	c.closeOnce.Do(
		func() {
			close(c.send)
		},
	)
}

func (c *Client) closeConnection() {
	if c.connection == nil {
		return
	}

	if err := c.connection.Close(); err != nil {
		log.Printf(
			"failed to close realtime connection: user_id=%d company_id=%d error=%v",
			c.userID,
			c.companyID,
			err,
		)
	}
}
