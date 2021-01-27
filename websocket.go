package kaiheila

import (
	"compress/zlib"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	"log"
	"time"
)

const (
	status_init = iota
	status_connecting
	status_connected
	status_retry
	status_close
	timeout = 6 * time.Second
	clock   = 100 * time.Millisecond
)

type HandlerFunc func(event EventMsg)

type WebSocketSession struct {
	conn *websocket.Conn

	api     *Client
	handler HandlerFunc

	sn       int
	status   int
	lastPong int64
}

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
		case status_init:
			if wss.conn != nil {
				wss.conn.Close()
			}
			wss.connect()
		case status_close:
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

	wss.status = status_connecting
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
		case SIG_EVENT:
			if wss.status == status_connected {
				wss.handler(msg.Data)
				wss.sn = msg.SN
			}
		case SIG_HELLO:
			if msg.Data.Code == 0 {
				wss.status = status_connected
				wss.lastPong = time.Now().Unix()
			} else {
				wss.status = status_init
				log.Println("[kaiheila] hello: %v", msg.Data.GetError())
			}
		case SIG_PONG:
			wss.status = status_connected
			wss.lastPong = time.Now().Unix()
		case SIG_RECONNECT:
			wss.status = status_init
		case SIG_RESUME_ACK:
			wss.status = status_connected
		}

		time.Sleep(clock)
	}
}

func (wss *WebSocketSession) ping() {
	for {
		time.Sleep(5 * timeout)

		// Check timeout
		switch wss.status {
		case status_init:
			return
		case status_connecting:
			wss.status = status_init
			log.Println("[kaiheila] no hello")
			return
		case status_connected:
			if wss.timeout() {
				wss.status = status_retry
				log.Println("[kaiheila] ping retry 1")
			}
		case status_retry:
			if wss.timeout() {
				wss.status = status_init
				log.Println("[kaiheila] ping retry 2")
				return
			}
		case status_close:
			return
		}

		err := wss.conn.WriteJSON(&websocketMsg{
			Signal: SIG_PING,
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

func (wss *WebSocketSession) Close() {
	wss.conn.Close()
	wss.status = status_close
}
