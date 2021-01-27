package kaiheila

import (
	"os"
	"testing"
	"time"
)

func TestKaiheila(t *testing.T) {
	client := NewClient("", TokenBot, os.Getenv("TEST_TOKEN"), 1)
	client.WebSocketSession(func(event EventMsg) {
		t.Log(event)
	})
	t.Log(client.SendChannelMsg(SendMessageReq{
		ChannelID: os.Getenv("TEST_CHANNEL"),
		Content:   "test",
	}))
	time.Sleep(10 * time.Second)
}
