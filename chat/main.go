package main

import (
	"fmt"

	"./utils"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gitlab.com/jinfagang/colorgo"
	"./db"
)

func Welcome() {
	cg.PrintlnGreen(cg.BoldStart + "Welcome to Uranus." + cg.BoldEnd)

	cg.PrintlnBlue(cg.BoldStart + `   __  __                           
  / / / /________ _____ __   ______
 / / / / ___/ __ / __ \/ / / / ___/
/ /_/ / /  / /_/ / / / / /_/ (__  ) 
\____/_/   \__,_/_/ /_/\__,_/____/  
									` + cg.BoldEnd)
	cg.PrintlnRed("Author by Lucas Jin.")
	fmt.Println()
	cg.ResetColor()
}

// setup things, init db connection, check system status etc.
func Setup(c *utils.AppConfig){
	log.Info("Setting up Uranus...")
	log.Info(c)
	r := db.ConnectDB(c)
	if !r {
		log.Panic("Database connect failed. Pls check, exit Uranus.")
	} else{
		log.Info("Success connected to database.")
	}
}

// build all routes here
func BuildRoutes() *gin.Engine{
	r := gin.Default()

	// set the v1 version group
	v1 := r.Group("/api/v1")
	{
		// user register and get
		v1.POST("/users", HandleUserRegister)
		v1.GET("/users", HandleUserGet)

		// user status and get
		v1.POST("/status", HandleUserRegister)
		v1.GET("/status", HandleUserGet)
	}

	// for future placeholder
	v2 := r.Group("/api/v2")
	{
		// user register and get
		v2.POST("/users", HandleUserRegister)
		v2.GET("/users", HandleUserGet)

		// user status and get
		v2.POST("/status", HandleUserRegister)
		v2.GET("/status", HandleUserGet)
	}

	// static file routes
	r.GET("/chat", ServeHome)

	// ws route
	r.GET("/v1/ws", HandleConnections)
	return r
}


func main() {
	gin.SetMode(gin.DebugMode)
	log.SetLevel(log.DebugLevel)

	Welcome()
	// get config and do setup things
	c := utils.Config("./config.toml")
	Setup(c)

	// APIs routes part
	route := BuildRoutes()
	log.Info(c.Server.Addr + c.Server.Port)

	// handle msgs
	go HandleMessages()

	// this command must call at last
	route.Run(c.Server.Port)
}
