package utils

import 	(
	log "github.com/sirupsen/logrus"
	"gitlab.com/jinfagang/colorgo"
)

func CheckError(err error, info string) {
	if err != nil {
		log.Fatalf("%sGot error%s: %s, detail: %s", cg.BoldStart, cg.BoldEnd, info, err.Error())
	}
}
