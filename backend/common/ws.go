package common

import (
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	clients       map[*websocket.Conn]bool
	clientMux     sync.RWMutex
	closeCh       chan struct{}
	lastHeartbeat map[*websocket.Conn]time.Time // 用于跟踪每个客户端的最后心跳时间
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		clients:       make(map[*websocket.Conn]bool),
		lastHeartbeat: make(map[*websocket.Conn]time.Time),
		closeCh:       make(chan struct{}),
	}
}

func (w *WebsocketService) Run() {
	logrus.Info("Websocket Service Start")
	http.HandleFunc("/ws", w.handleConnections)
	go w.cleanupTicker()
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		logrus.Errorf("ListenAndServe Error: %s", err)
	}
}

func (w *WebsocketService) handleConnections(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		logrus.Errorf("Upgrade Error: %s", err)
		return
	}
	defer ws.Close()

	// Register the new client
	w.clientMux.Lock()
	w.clients[ws] = true
	w.clientMux.Unlock()

	logrus.Info("New client connected")

	for {
		time.Sleep(1 * time.Second)

		var msg string
		_, p, err := ws.ReadMessage()
		if err != nil {
			logrus.Errorf("Websocket Read Error: %s", err.Error())
			w.clientMux.Lock()
			delete(w.clients, ws)
			delete(w.lastHeartbeat, ws)
			w.clientMux.Unlock()
			break
		}
		msg = string(p)

		logrus.Info("Received message: " + msg)

		if msg == "heartbeat" {
			w.clientMux.Lock()
			w.lastHeartbeat[ws] = time.Now()
			w.clientMux.Unlock()
		} else {
			w.BroadcastToClients(msg)
		}
	}
}

func (w *WebsocketService) BroadcastToClients(message string) {

	logrus.Infof("BroadcastToClients: %s", message)

	// Use RLock for reading the map.
	w.clientMux.RLock()
	clientsCopy := make([]*websocket.Conn, 0, len(w.clients))
	for client := range w.clients {
		clientsCopy = append(clientsCopy, client)
	}
	w.clientMux.RUnlock()

	// Now iterate over the copied slice.
	for _, client := range clientsCopy {
		err := client.WriteJSON(message)
		if err != nil {
			logrus.Errorf("Write Error: %s", err)
			w.clientMux.Lock()
			delete(w.clients, client)
			w.clientMux.Unlock()
		}
	}
}

// cleanupTicker 函数来检查哪些客户端已超时
func (w *WebsocketService) cleanupTicker() {
	ticker := time.NewTicker(30 * time.Second)
	heartbeatTimeout := 1 * time.Minute // 设置为你想要的超时时间
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 检查超时的客户端
			now := time.Now()
			w.clientMux.Lock()
			for client, lastBeat := range w.lastHeartbeat {
				if now.Sub(lastBeat) > heartbeatTimeout {
					logrus.Warn("Client timeout, closing connection.")
					delete(w.clients, client)
					delete(w.lastHeartbeat, client)
					client.Close()
				}
			}
			w.clientMux.Unlock()
		case <-w.closeCh:
			return
		}
	}
}

func (w *WebsocketService) Close() {
	close(w.closeCh)
	for client := range w.clients {
		client.Close()
	}
}
