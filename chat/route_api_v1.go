package main

import (
	"github.com/gin-gonic/gin"
	"./db"
	"./models"
	"strconv"
	"net/http"
	log "github.com/sirupsen/logrus"
	"./utils"
	//"os/user"
)

// Handle all simple APIs
// How to connect to local database anyway?

func HandleUserRegister(context *gin.Context) {

	userAcc := context.PostForm("user_acc")
	userAvatar := context.PostForm("user_avatar")
	userNickName := context.PostForm("user_nick_name")
	userPassword := context.PostForm("user_password")
	userPhone := context.PostForm("user_phone")
	userEmail := context.PostForm("user_email")
	userGender := context.PostForm("user_gender")
	userSign := context.PostForm("user_sign")

	userType := context.PostForm("user_type")
	userTypeInt, _ := strconv.Atoi(userType)

	if userAcc == "" || userNickName == "" || userPassword == ""{
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "invalid",
			"code": http.StatusBadRequest,
			"msg": "user_acc, user_nick_name, user_password must not be none",
			"data": "",
		})
	}
	user := models.User{
		UserAcc:userAcc,
		UserAvatar:userAvatar,
		UserNickName:userNickName,
		UserPassword:userPassword,
		UserPhone:userPhone,
		UserEmail:userEmail,
		UserGender:userGender,
		UserSign:userSign,
		UserType:models.UserType(userTypeInt),
	}
	userTry := models.User{}
	if db.DB.Where("user_acc=?", userAcc).First(&userTry).RecordNotFound(){
		// user not found, create it
		db.DB.Create(&user)
		uAddr := utils.GenAddr(user.ID)
		user.UserAddr = "usr" + uAddr

		log.Infof("FUCK GenAddr: %s gened: %s", user.UserAddr, uAddr)
		db.DB.Save(&user)

		// should return a token to user, as well as login
		claims := make(map[string]interface{})
		claims["id"] = user.ID
		claims["msg"] = "hiding egg"
		claims["user_addr"] = user.UserAddr
		token, _ := utils.Encrypt(claims)
		log.Infof("Request new user: %s, it is new.", user)
		data := map[string]interface{}{"token": token, "id": user.ID, "user_addr": user.UserAddr}
		context.JSON(200, gin.H{
			"status": "success",
			"code": http.StatusOK,
			"msg": "user register succeed.",
			"data": data,
		})
	}else{
		log.Info("user exist.")
		context.JSON(200, gin.H{
			"status": "conflict",
			"code": http.StatusConflict,
			"msg": "user already exist.",
			"data": nil,
		})
	}
}

// get user information operation
// call user info should provide user_name
func HandleUserGet(context *gin.Context) {
	context.JSON(200, gin.H{
		"status": "success",
		"message": "you want get user info.",
	})
}

// For user login, returns a token, this only needs while first login
// or token expired purpose.
func HandleUserLogin(context *gin.Context)  {
	userAcc := context.PostForm("user_acc")
	userPassword := context.PostForm("user_password")

	// find the user and check the password
	// if right, return a token, otherwise refuse login
	userTry := models.User{}
	if db.DB.Where("user_acc = ?", userAcc).First(&userTry).RecordNotFound(){
		context.JSON(200, gin.H{
			"status": "error",
			"code": http.StatusNotFound,
			"msg": "login what? you are not even exist!",
			"data": "",
		})
	} else {
		//log.Infof("[login] here what found: %s", userTry)
		if userTry.UserPassword == userPassword{
			// return a token?
			claims := make(map[string]interface{})
			claims["id"] = userTry.ID
			claims["msg"] = "hiding egg"
			claims["user_addr"] = userTry.UserAddr
			token, _ := utils.Encrypt(claims)
			data := map[string]interface{}{"token": token, "id": userTry.ID, "user_addr": userTry.UserAddr}
			context.JSON(200, gin.H{
				"status": "success",
				"code": http.StatusOK,
				"msg": "login success, welcome " + userTry.UserNickName,
				"data": data,
			})
		} else {
			// login failed, refuse it
			context.JSON(200, gin.H{
				"status": "unauthorized",
				"code": http.StatusUnauthorized,
				"msg": "you are not allowed to login",
				"data": "",
			})
		}
	}
}


// User edit operation
func HandleUserEdit(context *gin.Context) {

}
