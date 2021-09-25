package kaiheila

import (
	"log"
	"os"
	"testing"
)

func TestKaiheila(t *testing.T) {
	client := NewClient("", TokenBot, os.Getenv("TEST_TOKEN"), 1)
	client.WebSocketSession(func(event EventMsg) {
		log.Println(event)
		if event.ChannelType == ChannelPerson && event.Type == MsgTypeText {
			var res map[string]interface{}
			log.Println(client.request("POST", 3, "direct-message/create", &struct {
				Content  string `json:"content"`
				TargetID string `json:"target_id"`
			}{
				Content:  "Hello!",
				TargetID: event.AuthorID,
			}, &res))
			log.Println(res)
		}
		if event.Extra.Is(ExtraPrivateAddedReaction) {
			var res map[string]interface{}
			event.Extra.GetBody(&res)
			log.Println(res)
		}
	})
	select {}
}
