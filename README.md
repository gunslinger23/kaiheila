# kaiheila Golang SDK

⚠️仍在开发阶段⚠️

## Usage

```golang
// Init new client for kaiheila
client := NewClient("", TokenBot, "TOKEN", 1)

// Listen WebSocket event
client.WebSocketSession(func(event EventMsg) {
	fmt.Println(event)
    // Get extra
    if event.Extra.Is(MsgTypeText) {
        fmt.Println(event.Extra.Author.Username, ":", event.Content)
    }
    if event.Extra.Is(ExtraGuildMemberOnline) {
        fmt.Println(event.Extra.Body["user_id"], "is online!")
    }
    if event.Extra.Is(ExtraGuildMemberOffline) {
        fmt.Println(event.Extra.Body["user_id"], "is offline!")
    }
})

// Use http api
fmt.Println(client.SendChannelMsg(SendMessageReq{
	ChannelID: "ChannelID",
	Content:   "Hello world!",
}))
// Or use client request (for missing api)
// Structure must be tagged json field name
req := SendMessageReq{
	ChannelID: "ChannelID",
	Content:   "Hello world!",
}
resp := SendMessageResp{}
fmt.Println(c.request("POST", 3, "channel/message", &req, &resp))
fmt.Println(resp)

// Keep WebSocket goroutine running
time.Sleep(time.Minute)

client.Close()
```
