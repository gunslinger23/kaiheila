package kaiheila

// SendMessageReq message request struct
type SendMessageReq struct {
	Type         int    `json:"type,omitempty"`           // (Optional) A message type to send
	ChannelID    string `json:"channel_id"`               // Send to which channel?
	Content      string `json:"content"`                  // Content
	Quote        string `json:"quote,omitempty"`          // (Optional) Reply a message (msgID)
	Nonce        string `json:"nonce,omitempty"`          // (Optional) Server do not process message
	TempTargetID string `json:"temp_target_id,omitempty"` // (Optional) User id, the message will not store in server
}

// SendMessageResp message respone struct
type SendMessageResp struct {
	MsgID        string `json:"msg_id"`
	MsgTimestamp int64  `json:"msg_timestamp"`
	Nonce        string `json:"nonce"`
}

// SendChannelMsg Send a message to channel
func (c *Client) SendChannelMsg(req SendMessageReq) (SendMessageResp, error) {
	resp := &SendMessageResp{}
	err := c.request("POST", 3, "channel/message", struct2values(&req), resp)
	return *resp, err
}
