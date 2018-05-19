package utils

import (

	"github.com/dgrijalva/jwt-go"
	"errors"
	"time"
)

var mySalt = []byte("so fucking awesome")

// expiring time, 90 days
var exp = time.Hour * 24 * 90
// my domain
var iss = "loliloli.pro"

// Error handling
var TokenInvalid = errors.New("token invalid")
var TokenExpired = errors.New("token has been expired")
var TokenIllegal = errors.New("you are not allowed access our token")


// Claims is the content should to be encrypt
func Encrypt(infoMap map[string]interface{}) (string, error){
	/*
	claims are like this:

	claims := &jwt.StandardClaims{
    ExpiresAt: 15000,
    Issuer:    "test",}
	*/
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["iss"] = iss
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(exp).Unix()
	for k, v := range infoMap{
		claims[k] = v
	}
	token.Claims = claims
	ss, err := token.SignedString(mySalt)
	return ss, err
}

func Decrypt(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, getValidationKey)
	if err != nil {
		return nil, TokenInvalid
	}

	if jwt.SigningMethodHS256.Alg() != token.Header["alg"] {
		return nil, TokenInvalid
	}
	if !token.Valid {
		return nil, TokenInvalid
	}

	claims := token.Claims.(jwt.MapClaims)
	if claims["iss"] != iss {
		return nil, TokenIllegal
	}
	return claims, nil
}

func RefreshToken(token *jwt.Token) {
	// not implement
}

func getValidationKey(*jwt.Token) (interface{}, error) {
	return mySalt, nil
}
