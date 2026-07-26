package realtime

import (
	"encoding/json"
	"log"
	"sync"
)

type Hub struct {
	mu sync.RWMutex

	clientsByCompany map[uint64]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clientsByCompany: make(
			map[uint64]map[*Client]struct{},
		),
	}
}

func (h *Hub) Register(
	client *Client,
) {
	if client == nil ||
		client.companyID == 0 {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	companyClients, exists :=
		h.clientsByCompany[client.companyID]
	if !exists {
		companyClients = make(
			map[*Client]struct{},
		)

		h.clientsByCompany[client.companyID] =
			companyClients
	}

	companyClients[client] = struct{}{}

	log.Printf(
		"realtime client registered: user_id=%d company_id=%d company_connections=%d",
		client.userID,
		client.companyID,
		len(companyClients),
	)
}

func (h *Hub) Unregister(
	client *Client,
) {
	if client == nil ||
		client.companyID == 0 {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	companyClients, exists :=
		h.clientsByCompany[client.companyID]
	if !exists {
		return
	}

	if _, exists := companyClients[client]; !exists {
		return
	}

	delete(
		companyClients,
		client,
	)

	client.closeSend()

	if len(companyClients) == 0 {
		delete(
			h.clientsByCompany,
			client.companyID,
		)
	}

	log.Printf(
		"realtime client unregistered: user_id=%d company_id=%d company_connections=%d",
		client.userID,
		client.companyID,
		len(companyClients),
	)
}

func (h *Hub) BroadcastToCompany(
	companyID uint64,
	event Event,
) error {
	if companyID == 0 {
		return nil
	}

	message, err := json.Marshal(
		event,
	)
	if err != nil {
		return err
	}

	clients := h.companyClientsSnapshot(
		companyID,
	)

	for _, client := range clients {
		if !client.enqueue(message) {
			h.Unregister(client)
		}
	}

	return nil
}

func (h *Hub) BroadcastWorkContextChanged(
	companyID uint64,
	payload WorkContextChangedPayload,
) error {
	return h.BroadcastToCompany(
		companyID,
		Event{
			Type:    EventWorkContextChanged,
			Payload: payload,
		},
	)
}

func (h *Hub) BroadcastPresenceChanged(
	companyID uint64,
	payload PresenceChangedPayload,
) error {
	return h.BroadcastToCompany(
		companyID,
		Event{
			Type:    EventPresenceChanged,
			Payload: payload,
		},
	)
}

func (h *Hub) ConnectionCount(
	companyID uint64,
) int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(
		h.clientsByCompany[companyID],
	)
}

func (h *Hub) companyClientsSnapshot(
	companyID uint64,
) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	companyClients, exists :=
		h.clientsByCompany[companyID]
	if !exists {
		return nil
	}

	clients := make(
		[]*Client,
		0,
		len(companyClients),
	)

	for client := range companyClients {
		clients = append(
			clients,
			client,
		)
	}

	return clients
}
