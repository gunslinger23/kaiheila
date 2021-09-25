package kaiheila

import (
	"compress/zlib"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	statusInit = iota
	statusConnecting
	statusConnected
	statusRetry
	statusClose
	timeout = 6 * time.Second        // Timeout of connection
	clock   = 100 * time.Millisecond // Session status check clock
)

// HandlerFunc handle event from server
type HandlerFunc func(event EventMsg)

// WebSocketSession a websocket session
type WebSocketSession struct {
	wsCloser func() error

	api     *Client
	handler HandlerFunc

	sn       int
	status   int
	lastPong int64
}

// WebSocketSession Create a websocket session for handle event from server
func (c *Client) WebSocketSession(handler HandlerFunc) *WebSocketSession {
	wss := &WebSocketSession{
		api:     c,
		handler: handler,
	}
	go wss.run()
	return wss
}

func (wss *WebSocketSession) run() {
	for {
		switch wss.status {
		case statusInit:
			if wss.wsCloser != nil {
				wss.wsCloser()
			}
			// keep trying...
			for !wss.connect() && wss.status == statusInit {
				time.Sleep(timeout)
			}
		case statusClose:
			return
		}
		time.Sleep(clock)
	}
}

func (wss *WebSocketSession) connect() bool {
	gateway, err := wss.api.GetGateway()
	if err != nil {
		log.Println("[kaiheila] gateway:", err)
		return false
	}
	conn, _, err := websocket.DefaultDialer.Dial(gateway, nil)
	if err != nil {
		log.Println("[kaiheila] dial:", err)
		return false
	}

	wss.wsCloser = conn.Close
	wss.status = statusConnecting

	go wss.receive(conn)
	go wss.healthChecker(conn)
	return true
}

func (wss *WebSocketSession) receive(conn *websocket.Conn) {
	defer func() {
		wss.status = statusInit
	}()

	for {
		msg := &websocketMsg{}
		_, r, err := conn.NextReader()
		if err != nil {
			log.Println("[kaiheila] read:", err)
			return
		}
		raw, err := zlib.NewReader(r)
		if err != nil {
			log.Println("[kaiheila] zlib:", err)
		}
		err = json.NewDecoder(raw).Decode(msg)
		if err != nil {
			log.Println("[kaiheila] json:", err)
		}

		// signal
		switch msg.Signal {
		case signalEvent:
			if wss.status == statusConnected {
				wss.handler(msg.Data)
				wss.sn = msg.SN
			}
		case signalHello:
			if msg.Data.Code == 0 {
				wss.status = statusConnected
			} else {
				log.Println("[kaiheila] hello:", msg.Data.GetError())
				return
			}
		case signalPong:
			wss.status = statusConnected
			wss.lastPong = time.Now().Unix()
		case signalReconnect:
			return
		}
	}
}

func (wss *WebSocketSession) healthChecker(conn *websocket.Conn) {
	defer func() {
		wss.status = statusInit
	}()

	// Check hello
	time.Sleep(timeout)
	if wss.status != statusConnected {
		return
	}

	// Check ping pong
	for {
		err := conn.WriteJSON(&websocketMsg{
			Signal: signalPing,
			SN:     wss.sn,
		})
		if err != nil {
			log.Println("[kaiheila] ping:", err)
			return
		}

		time.Sleep(timeout)

		// Check timeout
		switch wss.status {
		case statusInit, statusClose, statusConnecting:
			return
		case statusConnected:
			if wss.timeout() {
				wss.status = statusRetry
				log.Println("[kaiheila] ping retry 1")
			}
		case statusRetry:
			if wss.timeout() {
				log.Println("[kaiheila] ping retry 2")
				return
			}
		}

		// sleep for next ping
		time.Sleep(4 * timeout)
	}
}

func (wss *WebSocketSession) timeout() bool {
	return wss.lastPong < time.Now().Add(-timeout).Unix()
}

// Close close this seesion
func (wss *WebSocketSession) Close() {
	if wss.wsCloser != nil {
		wss.wsCloser()
	}
	wss.status = statusClose
}
