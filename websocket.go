package kaiheila

import (
	"compress/zlib"
	"log"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

const (
	statusInit = iota
	statusConnecting
	statusConnected
	statusRetry
	statusClose
	timeout = 6 * time.Second
	clock   = 100 * time.Millisecond
)

// HandlerFunc handle event from server
type HandlerFunc func(event EventMsg)

// WebSocketSession a websocket session
type WebSocketSession struct {
	conn *websocket.Conn

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
			if wss.conn != nil {
				wss.conn.Close()
			}
			wss.connect()
		case statusClose:
			return
		}
		time.Sleep(clock)
	}
}

func (wss *WebSocketSession) connect() {
	gateway, err := wss.api.GetGateway()
	if err != nil {
		log.Println("[kaiheila] gateway:", err)
		return
	}
	wss.conn, _, err = websocket.DefaultDialer.Dial(gateway, nil)
	if err != nil {
		log.Println("[kaiheila] dial:", err)
		return
	}

	wss.status = statusConnecting
	go wss.receive()
	go wss.ping()
}

func (wss *WebSocketSession) receive() {
	for {
		msg := &websocketMsg{}
		_, r, err := wss.conn.NextReader()
		if err != nil {
			log.Println("[kaiheila] read:", err)
			return
		}
		raw, err := zlib.NewReader(r)
		if err != nil {
			log.Println("[kaiheila] zlib:", err)
		}
		err = jsoniter.NewDecoder(raw).Decode(msg)
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
				wss.lastPong = time.Now().Unix()
			} else {
				wss.status = statusInit
				log.Println("[kaiheila] hello:", msg.Data.GetError())
			}
		case signalPong:
			wss.status = statusConnected
			wss.lastPong = time.Now().Unix()
		case signalReconnect:
			wss.status = statusInit
		case signalResumeACK:
			wss.status = statusConnected
		}

		time.Sleep(clock)
	}
}

func (wss *WebSocketSession) ping() {
	for {
		time.Sleep(5 * timeout)

		// Check timeout
		switch wss.status {
		case statusInit:
			return
		case statusConnecting:
			wss.status = statusInit
			log.Println("[kaiheila] no hello")
			return
		case statusConnected:
			if wss.timeout() {
				wss.status = statusRetry
				log.Println("[kaiheila] ping retry 1")
			}
		case statusRetry:
			if wss.timeout() {
				wss.status = statusInit
				log.Println("[kaiheila] ping retry 2")
				return
			}
		case statusClose:
			return
		}

		err := wss.conn.WriteJSON(&websocketMsg{
			Signal: signalPing,
			SN:     wss.sn,
		})
		if err != nil {
			log.Println("[kaiheila] ping:", err)
		}
	}
}

func (wss *WebSocketSession) timeout() bool {
	return wss.lastPong < time.Now().Add(-5*timeout).Unix()
}

// Close close this seesion
func (wss *WebSocketSession) Close() {
	wss.conn.Close()
	wss.status = statusClose
}
