package models

import (
	"github.com/jinzhu/gorm"
	"time"

)

// the group model, contains all groups information
type GroupType int
const (
	// a feeds group only send feeds by robots
	GroupFeeds GroupType = iota
	// normal group
	GroupNormal
)

type Group struct{
	gorm.Model
	GroupName string `json:"group_name"`
	GroupCreateTime time.Time `json:"group_create_time"`
	GroupMembers []User `json:"group_members"`
	GroupAddr string `json:"group_addr"`
}
