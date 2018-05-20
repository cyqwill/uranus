package utils

import (

	"math/rand"
	"time"
)


var allChars = []string{
	"2", "3", "4", "5", "6", "7", "8", "9",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "l", "s", "t", "u",
	"v", "w", "x", "y", "z",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "L", "S", "T", "U",
	"V", "W", "X", "Y", "Z"}

/* generate UserAddr with UserId */
func GenAddr(id uint) string{
	// algorithm is simple, we using 2-9 + a-z + A-Z as 60 decimal
	// 8 length Address
	// 23 -> usr8fugkitU
	resultStr := ""
	for id / 60 > 0{
		tmpMod := id % 60
		resultStr += allChars[tmpMod]
		id /= 60
	}
	resultStr += randomStr(10 - len(resultStr))
	return resultStr
}

func randomStr(n int) string{
	rand.Seed(time.Now().UnixNano())
	r := ""
	i := 0
	for i < n {
		r += allChars[rand.Intn(len(allChars))]
		i += 1
	}
	return r
}
