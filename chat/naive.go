package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"./utils"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"time"
	"./models"
	"github.com/mitchellh/mapstructure"
)

// naive for test
// very simple websocket msg send implementation

type Client struct{
	Conn *websocket.Conn
	ClientAddr string
	IsOnline bool
	UA string
	// client online time
	UpTime int64
	LastSeenTime int64
}

// save current online clients
//var clients = make(map[*websocket.Conn]bool)
var clientsPool = make(map[string]*Client)
// msg queue are all the messages
var msgQueue = make(chan OutMsg)
var wsUpGrader = websocket.Upgrader{}


type IncomeMsg struct{
	// could be `send`, `hi`, `del`, `add`
	MsgType string `json:"msg_type"`
	Payload map[string]interface{} `json:"payload"`
}
type OutMsg struct{
	TargetAddr string
	MsgType int
	Payload interface{}
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

	// waite forever to write ws message to broadcast
	// How to let it support for group chat???
	for{
		var msg IncomeMsg
		// every incoming msg must provide a token and a target id (this 2 are very important)
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Errorf("Error in read websocket json, error: %s", err.Error())
			//delete(clients, ws)
			break
		} else {
			payload := msg.Payload
			switch msg.MsgType {
			case "send":
				// send msg,
				log.Infof("[send] incoming>: %s", payload)
				sendMsg := models.SendMsg{}
				err := mapstructure.Decode(payload, &sendMsg)
				if err != nil {
					// send back to same address a error msg say, hi msg not right
					var msg = OutMsg{}
					msg.TargetAddr = sendMsg.SendAddr
					msg.Payload = map[string]string{
						"error": "send msg error " + err.Error(),
					}
					msgQueue <- msg
				} else {
					// TODO: I am stuck here, how to send off-line msg?
					var msg = OutMsg{}
					msg.TargetAddr = sendMsg.TargetAddr
					msg.Payload = sendMsg.Content
					msgQueue <- msg
				}

			case "hi":
				log.Infof("[hi] incoming>: %s", payload)
				// add this to clients pool
				hiMsg := models.HiMsg{}
				err := mapstructure.Decode(payload, &hiMsg)
				if err != nil {
					// send back to same address a error msg say, hi msg not right
					var msg = OutMsg{}
					msg.TargetAddr = hiMsg.UserAddr
					msg.Payload = map[string]string{
						"error": "hi msg error " + err.Error(),
					}
					msgQueue <- msg
				}else{
					// new clients online
					claims, err := utils.Decrypt(hiMsg.Token)
					if err != nil {
						// token invalid, expired or not right, refuse connection
						var msg = OutMsg{}
						msg.TargetAddr = hiMsg.UserAddr
						msg.Payload = map[string]string{
							"error": "token invalid. " + err.Error(),
						}
						msgQueue <- msg
					} else{
						clientAddr := claims["user_addr"].(string)
						ua := hiMsg.UA
						client := Client{Conn: ws, ClientAddr:clientAddr, UA:ua, UpTime:time.Now().Unix()}
						clientsPool[clientAddr] = &client
					}
				}
			case "add":
				log.Infof("[add] incoming>: %s", payload)
				continue
			case "del":
				log.Infof("[del] incoming>: %s", payload)
				continue
			}
		}
	}
}

func HandleMessages(){
	for {
		//msg := <- broadcast
		msg := <- msgQueue
		log.Info("msg from msgQueue: ", msg)
		// TODO: What if targetAddr not in clientsPool????
		if targetClient, ok := clientsPool[msg.TargetAddr]; ok {
			err := targetClient.Conn.WriteJSON(msg.Payload)
			utils.CheckError(err, "msgQueue write json.")
		} else {
			// targetAddr not in clientsPool
			// I think when it on-line it will receive msg again!!
		}
	}
}


// Find client by address from clientsPool when it is online
// else return nil, indicates the client are off-line
func findClientFromPool(targetAddr string) {

}

