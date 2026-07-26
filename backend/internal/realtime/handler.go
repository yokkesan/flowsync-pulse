package realtime

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"flowsync-pulse/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	connectionTicketByteLength = 32
	connectionTicketValidTime  = 60 * time.Second
)

type Handler struct {
	hub         *Hub
	ticketStore ConnectionTicketStore
	upgrader    websocket.Upgrader
}

type CreateConnectionTicketResponse struct {
	Ticket    string    `json:"ticket"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewHandler(
	hub *Hub,
	ticketStore ConnectionTicketStore,
	allowedOrigin string,
) *Handler {
	return &Handler{
		hub:         hub,
		ticketStore: ticketStore,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(
				request *http.Request,
			) bool {
				return isAllowedWebSocketOrigin(
					request,
					allowedOrigin,
				)
			},
		},
	}
}

// CreateConnectionTicket godoc
// @Summary WebSocket接続チケット発行
// @Description 認証済みユーザー向けに、短時間かつ1回限り有効なWebSocket接続チケットを発行します。
// @Tags realtime
// @Produce json
// @Security BearerAuth
// @Success 201 {object} CreateConnectionTicketResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/realtime/tickets [post]
func (h *Handler) CreateConnectionTicket(
	c *gin.Context,
) {
	userID := c.GetUint64(
		middleware.ContextUserIDKey,
	)
	companyID := c.GetUint64(
		middleware.ContextCompanyIDKey,
	)

	if userID == 0 || companyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "ログイン情報が無効です。",
			},
		)
		return
	}

	token, err := generateConnectionTicketToken()
	if err != nil {
		log.Printf(
			"failed to generate realtime connection ticket: user_id=%d company_id=%d error=%v",
			userID,
			companyID,
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "接続チケットの発行に失敗しました。",
			},
		)
		return
	}

	now := time.Now().UTC()
	expiresAt := now.Add(
		connectionTicketValidTime,
	)

	if err := h.ticketStore.Create(
		c.Request.Context(),
		ConnectionTicket{
			Token:     token,
			UserID:    userID,
			CompanyID: companyID,
			ExpiresAt: expiresAt,
			UsedAt:    nil,
		},
	); err != nil {
		log.Printf(
			"failed to save realtime connection ticket: user_id=%d company_id=%d error=%v",
			userID,
			companyID,
			err,
		)

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"message": "接続チケットの発行に失敗しました。",
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated,
		CreateConnectionTicketResponse{
			Ticket:    token,
			ExpiresAt: expiresAt,
		},
	)
}

// Connect godoc
// @Summary WebSocket接続
// @Description 発行済み接続チケットを使用してWebSocketへ接続します。
// @Tags realtime
// @Param ticket query string true "WebSocket接続チケット"
// @Success 101
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/realtime/ws [get]
func (h *Handler) Connect(
	c *gin.Context,
) {
	ticketToken := strings.TrimSpace(
		c.Query("ticket"),
	)
	if ticketToken == "" {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"message": "接続チケットが必要です。",
			},
		)
		return
	}

	now := time.Now().UTC()

	ticket, err := h.ticketStore.Consume(
		c.Request.Context(),
		ticketToken,
		now,
	)
	if err != nil {
		log.Printf(
			"realtime connection ticket rejected: client_ip=%s error=%v",
			c.ClientIP(),
			err,
		)

		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "接続チケットが無効です。",
			},
		)
		return
	}

	if ticket.UserID == 0 ||
		ticket.CompanyID == 0 {
		c.JSON(
			http.StatusUnauthorized,
			gin.H{
				"message": "接続チケットが無効です。",
			},
		)
		return
	}

	connection, err := h.upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)
	if err != nil {
		log.Printf(
			"failed to upgrade realtime connection: user_id=%d company_id=%d client_ip=%s error=%v",
			ticket.UserID,
			ticket.CompanyID,
			c.ClientIP(),
			err,
		)
		return
	}

	client := NewClient(
		h.hub,
		connection,
		ticket.UserID,
		ticket.CompanyID,
	)

	client.Start()
}

func generateConnectionTicketToken() (
	string,
	error,
) {
	randomBytes := make(
		[]byte,
		connectionTicketByteLength,
	)

	if _, err := rand.Read(
		randomBytes,
	); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(
		randomBytes,
	), nil
}

func isAllowedWebSocketOrigin(
	request *http.Request,
	allowedOrigin string,
) bool {
	origin := strings.TrimSpace(
		request.Header.Get("Origin"),
	)

	// curl等のOriginを送らないクライアントは拒否せず、
	// 接続チケットによる認証を必須とする。
	if origin == "" {
		return true
	}

	normalizedAllowedOrigin := strings.TrimRight(
		strings.TrimSpace(allowedOrigin),
		"/",
	)
	normalizedOrigin := strings.TrimRight(
		origin,
		"/",
	)

	if normalizedAllowedOrigin == "" {
		return false
	}

	return normalizedOrigin ==
		normalizedAllowedOrigin
}
