package models

import "github.com/jinzhu/gorm"

// user friends struct

type Friends struct{
	gorm.Model
	UserAId int
	UserBId int
	IsFriends bool
}