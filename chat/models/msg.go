package models

import (
	_ "encoding/json"
	"github.com/jinzhu/gorm"
)

type SendType int
type MsgType int

const (
	// one to one
	O2O SendType = iota
	// one to many, for group chat
	O2M
	// many to one, for official account receive
	M2O
)

const (
	Text MsgType = iota
	Image
	Video
	UrlLink
	Contact
	System
	Notification
	Sticker
	// Choice button for quicker answer
	ChoiceButton
	// Command order
	Command
)

// normal msg struct for db storage


type Msg struct{
	gorm.Model
	SenderId   int      `json:"sender_id"`
	TargetId   int      `json:"target_id"`
	SendTime   string   `json:"send_time"`
	ReadTime   string   `json:"read_time"`
	MsgContent string   `json:"msg_content"`
	SendType    SendType `json:"send_type"`
	MsgType int `json:"msg_type"`
}

type TextMsg struct{
	MsgId      int      `json:"msg_id"`
	SenderId   int      `json:"sender_id"`
	SendTime   string   `json:"send_time"`
	ReadTime   string   `json:"read_time"`
	MsgContent string   `json:"msg_content"`
	SendType    SendType `json:"send_type"`
	MsgType int `json:"msg_type"`
}

type ImageMsg struct {
	MsgId      int      `json:"msg_id"`
	SenderId   int      `json:"sender_id"`
	SendTime   string   `json:"send_time"`
	ReadTime   string   `json:"read_time"`
	MsgContent byte     `json:"msg_content"`
	MsgSize    int      `json:"msg_size"`
	SendType    SendType `json:"send_type"`
	MsgType int `json:"msg_type"`
}

// HiMsg for token validation
// HiMsg is very important, cause I have to know which/kind devices it is
// Phone or PC or Arduino or RaspberryPi, even the location should provide
// UA include the operation system info and uranus client version
type HiMsg struct{
	Token string `json:"token"`
	UserAddr string `json:"user_addr"`
	Device string `json:"device"`
	Location string `json:"location"`
	UA string `json:"ua"`
}

// Generous sending msg
type SendMsg struct{
	Target string `json:"target"`
	Sender string `json:"sender"`
	Content string `json:"content"`
	MsgType int `json:"msg_type"`
}

// For adding a group or something
type AddMsg struct{

}

// For deleting msg
type DelMsg struct{

}
