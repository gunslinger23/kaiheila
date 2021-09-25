package kaiheila

import (
	"encoding/json"
	"fmt"
)

const (
	signalEvent = iota
	signalHello
	signalPing
	signalPong
	signalResume
	signalReconnect
	signalResumeACK
)

type websocketMsg struct {
	Signal int      `json:"s"`
	Data   EventMsg `json:"d"`
	SN     int      `json:"sn"`
}

type httpMsg struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

// MsgType message type
type MsgType int

const (
	// MsgTypeText Type: Text
	MsgTypeText = MsgType(1)
	// MsgTypeImg Type: Image
	MsgTypeImg = MsgType(2)
	// MsgTypeVideo Type: Video
	MsgTypeVideo = MsgType(3)
	// MsgTypeFile Type: File
	MsgTypeFile = MsgType(4)
	// MsgTypeVoice Type: Voice
	MsgTypeVoice = MsgType(8)
	// MsgTypeKmarkdown Type: Kmarkdown
	MsgTypeKmarkdown = MsgType(9)
	// MsgTypeCard Type: Card
	MsgTypeCard = MsgType(10)
	// MsgTypeSystem Type: System
	MsgTypeSystem = MsgType(255)
)

// ChannelType Channel type
type ChannelType string

const (
	ChannelGroup  = ChannelType("GROUP")
	ChannelPerson = ChannelType("PERSON")
)

// EventMsg Event message from server
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

// ExtraType non-system extra type
type ExtraType string

const (
	ExtraAddedReaction          = ExtraType("added_reaction")
	ExtraDeletedReaction        = ExtraType("deleted_reaction")
	ExtraUpdateMessage          = ExtraType("updated_message")
	ExtraDeletedMessage         = ExtraType("deleted_message")
	ExtraAddedChannel           = ExtraType("added_channel")
	ExtraUpdatedChannel         = ExtraType("updated_channel")
	ExtraDeletedChannel         = ExtraType("deleted_channel")
	ExtraPinnedMessage          = ExtraType("pinned_message")
	ExtraUnpinnedMessage        = ExtraType("unpinned_message")
	ExtraUpdatePrivateMessage   = ExtraType("updated_private_message")
	ExtraDeletedPrivateMessage  = ExtraType("deleted_private_message")
	ExtraPrivateAddedReaction   = ExtraType("private_added_reaction")
	ExtraPrivateDeletedReaction = ExtraType("private_deleted_reaction")
	ExtraJoinedGuild            = ExtraType("joined_guild")
	ExtraExitedGuild            = ExtraType("exited_guild")
	ExtraUpdateGuildMember      = ExtraType("updated_guild_member")
	ExtraGuildMemberOnline      = ExtraType("guild_member_online")
	ExtraGuildMemberOffline     = ExtraType("guild_member_offline")
	ExtraAddedRole              = ExtraType("added_role")
	ExtraDeletedRole            = ExtraType("deleted_role")
	ExtraUpdatedRole            = ExtraType("updated_role")
	ExtraUpdatedGuild           = ExtraType("updated_guild")
	ExtraDeletedGuild           = ExtraType("deleted_guild")
	ExtraAddedBlockList         = ExtraType("added_block_list")
	ExtraDeletedBlockList       = ExtraType("deleted_block_list")
	ExtraJoinedChannel          = ExtraType("joined_channel")
	ExtraExitedChannel          = ExtraType("exited_channel")
	ExtraUserUpdated            = ExtraType("user_updated")
	ExtraSelfJoinedGuild        = ExtraType("self_joined_guild")
	ExtraSelfExitedGuild        = ExtraType("self_exited_guild")
	ExtraMessageButtonClick     = ExtraType("message_btn_click")
)

// ExtraMsg extra info of message
type ExtraMsg struct {
	Type         json.RawMessage `json:"type"`
	Body         json.RawMessage `json:"body"`
	GuildID      string          `json:"guild_id"`
	ChannelName  string          `json:"channel_name"`
	Mention      []string        `json:"mention"`
	MentionAll   bool            `json:"mention_all"`
	MentionRoles []string        `json:"mention_roles"`
	MentionHere  bool            `json:"mention_here"`
	Author       AuthorMsg       `json:"author"`
}

// Is match type of non-system message
func (msg ExtraMsg) Is(et ExtraType) bool {
	var res string
	_ = json.Unmarshal(msg.Type, &res)
	return res == string(et)
}

// GetBody get body of extra message
func (msg ExtraMsg) GetBody(dest interface{}) error {
	return json.Unmarshal(msg.Body, dest)
}

// AuthorMsg Author info of message
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

// GetError Get error from message
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
