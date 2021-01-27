package kaiheila

import "fmt"

const (
	SIG_EVENT = iota
	SIG_HELLO
	SIG_PING
	SIG_PONG
	SIG_RESUME
	SIG_RECONNECT
	SIG_RESUME_ACK
)

type websocketMsg struct {
	Signal int      `json:"s"`
	Data   EventMsg `json:"d"`
	SN     int      `json:"sn"`
}

type httpMsg struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type MsgType int

const (
	MsgTypeText      = MsgType(1)
	MsgTypeImg       = MsgType(2)
	MsgTypeVideo     = MsgType(3)
	MsgTypeFile      = MsgType(4)
	MsgTypeVoice     = MsgType(8)
	MsgTypeKmarkdown = MsgType(9)
	MsgTypeSystem    = MsgType(255)
)

type ChannelType string

const (
	ChannelGroup = ChannelType("GROUP")
)

type EventMsg struct {
	// Signal
	Code      int    `json:"code"`
	SessionID string `json:"sessionId"`
	Error     string `json:"err"`
	// Server push
	ChannelType  ChannelType `json:"channel_type"`
	Type         MsgType     `json:"type"`
	TargetID     string      `json:"target_id"` // GROUP: channel_id
	AuthorID     string      `json:"author_id"`
	Content      string      `json:"content"`
	MsgID        string      `json:"msg_id"`
	MsgTimestamp int64       `json:"msg_timestamp"`
	Nonce        string      `json:"nonce"`
	Extra        ExtraMsg    `json:"extra"`
}

type ExtraMsg struct {
	Type         MsgType   `json:"type"`
	GuildID      string    `json:"guild_id"`
	ChannelName  string    `json:"channel_name"`
	Mention      []string  `json:"mention"`
	MentionAll   bool      `json:"mention_all"`
	MentionRoles []string  `json:"mention_roles"`
	MentionHere  bool      `json:"mention_here"`
	Author       AuthorMsg `json:"author"`
}

type AuthorMsg struct {
	ID          string   `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	IdentifyNum string   `json:"identify_num"`
	Online      bool     `json:"online"`
	Avatar      string   `json:"avatar"`
	Roles       []string `json:"roles"`
	Bot         bool     `json:"bot"`
}

func (msg EventMsg) GetError() error {
	switch msg.Code {
	case 40100:
		return fmt.Errorf("missing arg")
	case 40101:
		return fmt.Errorf("invalid token")
	case 40102:
		return fmt.Errorf("token auth failed")
	case 40103:
		return fmt.Errorf("token expired")
	case 40106:
		return fmt.Errorf("resume failed, missing arg(%s)", msg.Error)
	case 40107:
		return fmt.Errorf("resume failed, session expired(%s)", msg.Error)
	case 40108:
		return fmt.Errorf("resume failed, invalid sn(%s)", msg.Error)
	default:
		return nil
	}
}