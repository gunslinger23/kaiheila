package kaiheila

// SendMessageReq Message request struct
type SendMessageReq struct {
	Type         int    `json:"type,omitempty"`           // (Optional) A message type to send
	ChannelID    string `json:"channel_id"`               // Send to which channel?
	Content      string `json:"content"`                  // Content
	Quote        string `json:"quote,omitempty"`          // (Optional) Reply a message (msgID)
	Nonce        string `json:"nonce,omitempty"`          // (Optional) Server do not process message
	TempTargetID string `json:"temp_target_id,omitempty"` // (Optional) User id, the message will not store in server
}

// SendMessageResp Message respone struct
type SendMessageResp struct {
	MsgID        string `json:"msg_id"`        // ID of message sent
	MsgTimestamp int64  `json:"msg_timestamp"` // Timestamp of message sent
	Nonce        string `json:"nonce"`         // Server do not process message
}

// SendChannelMsg Send a message to channel
func (c *Client) SendChannelMsg(req SendMessageReq) (SendMessageResp, error) {
	resp := &SendMessageResp{}
	err := c.request("POST", 3, "channel/message", &req, resp)
	return *resp, err
}
