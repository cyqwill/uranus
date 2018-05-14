package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"./utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// naive for test
// very simple websocket msg send implementation

// save current online clients
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan TestMsg)
var wsUpGrader = websocket.Upgrader{}


type TestMsg struct{
	SenderName string `json:"sender_name"`
	Content string `json:"content"`
}

func ServeHome(c *gin.Context) {
	fmt.Println(c.Request.URL)
	if c.Request.Method != "GET" {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	//fs := http.FileServer(http.Dir("../public"))
	//http.Handle("/", fs)
	http.ServeFile(c.Writer, c.Request, "chat/public/index.html")
}

func HandleConnections(c *gin.Context){
	ws, err := wsUpGrader.Upgrade(c.Writer, c.Request, nil)
	utils.CheckError(err, "naive.HandleConnections")
	defer ws.Close()
	clients[ws] = true

	// waite forever to write ws message to broadcast
	for{
		var msg TestMsg
		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Printf("error in read ws json msg, format error: %v", err)
			//delete(clients, ws)
			break
		} else {
			log.Info("incoming>: ", msg)
		}
		broadcast <- msg
	}
}

func HandleMessages(){
	for {
		msg := <- broadcast
		log.Info("msg from broadcast: ", msg)
		for client := range clients{
			err := client.WriteJSON(msg)
			if err != nil{
				// if write error, then delete this client.
				fmt.Printf("error in write msg to client's ws: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}