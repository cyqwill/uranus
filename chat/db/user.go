package db

import (
	"log"
	"fmt"
)

func init() {
	log.SetPrefix("[main center]")
	log.SetFlags(log.Ldate | log.Lshortfile)
}

func main() {
	fmt.Println("hello, Go lang.")
}
