package main

import (
	"github.com/gin-gonic/gin"
	"./db"
	"./models"
	"strconv"
	"net/http"
	log "github.com/sirupsen/logrus"
)

// Handle all simple APIs
// How to connect to local database anyway?

func HandleUserRegister(context *gin.Context) {

	userAcc := context.PostForm("user_acc")
	userAvatar := context.PostForm("user_avatar")
	userNickName := context.PostForm("user_nick_name")
	userPhone := context.PostForm("user_phone")
	userEmail := context.PostForm("user_email")
	userGender := context.PostForm("user_gender")
	userSign := context.PostForm("user_sign")

	userType := context.PostForm("user_type")
	userTypeInt, _ := strconv.Atoi(userType)

	// gen user addr for the register

	if userAcc == "" || userNickName == ""{
		context.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusOK,
			"message": "user_acc and user_nick_name must not be none",
			"data": "",
		})
	} else {
		user := models.User{
			UserAcc:userAcc,
			UserAvatar:userAvatar,
			UserNickName:userNickName,
			UserPhone:userPhone,
			UserEmail:userEmail,
			UserGender:userGender,
			UserSign:userSign,
			UserType:models.UserType(userTypeInt),
		}
		// new an user
		log.Info(user)
		db.DB.Save(&user)

		// then return token for current user id
		data := map[string]interface{}{"token": "not ready yet", "id": user.UserNickName}
		context.JSON(200, gin.H{
			"status": http.StatusOK,
			"message": "user register succeed.",
			"data": data,
		})
	}

}

func HandleUserGet(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "success",
		"message": "you want get user info.",
	})
}
