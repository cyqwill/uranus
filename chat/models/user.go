package models

import (
	_ "encoding/json"
	"github.com/jinzhu/gorm"
)

type UserType int
const (
	Human UserType = iota
	// the most human father
	HumanFather
	// an robot father is the most powerful robot, it can send msg to any people without permission
	RobotFather
	Robot
)

func (u UserType) String() string{
	names := [...]string{
		"Human",
		"HumanFather",
		"RobotFather",
		"Robot",
	}
	return names[u]
}

type UserTags struct{
	Tag string
}


type User struct {
	gorm.Model
	// an user account is the only way to find him
	UserAcc string  `json:"user_acc"`

	// using base64 save user avatar
	UserAvatar string `json:"user_avatar"`
	UserPhone string `json:"user_phone"`
	UserEmail string `json:"user_email"`
	UserPassword string `json:"user_password"`
	// tags for finding people by tag
	//UserTags []string `json:"user_tags"`
	UserNickName string `json:"user_nick_name"`
	UserGender string `json:"user_gender"`
	UserSign string `json:"user_sign"`
	UserBirth string `json:"user_birth"`
	UserCity string `json:"user_city"`
	// user type indicates what user is this
	UserType UserType `json:"user_type"`
	// user addr simply a map of user_id, but with encode and hide while using http communicate
	UserAddr string `json:"user_addr"`
}


