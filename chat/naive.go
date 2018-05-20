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
	"io"
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
// Note: clientsPool now the value is a client list
var clientsPool = make(map[string][]*Client)
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

type TestSendMsg struct{
	Target string
	Sender string
	Content string
	MsgType int
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
			if err == io.ErrUnexpectedEOF{
				// indicates this connection are closed, should delete it from clientsPool
				ws.Close()
				log.Errorf("got a closed connection %s", ws)
			} else {
				log.Errorf("Error in read websocket json, error: %s", err.Error())
				errorResponse := map[string]string{
					"error": "your must send json as bytes to server through websocket.",
				}
				ws.WriteJSON(errorResponse)
				break
			}
		} else {
			payload := msg.Payload
			switch msg.MsgType {
			case "send":
				// send msg,
				log.Infof("[send] incoming>: %s", payload)
				sendMsg := models.SendMsg{}
				err := mapstructure.Decode(payload, &sendMsg)
				log.Infof("sendMsg: %s", sendMsg)
				if err != nil {
					// send back to same address a error msg say, hi msg not right
					errorResponse := map[string]string{
						"error": "your send msg struct is not right, read uranus doc first.",
					}
					ws.WriteJSON(errorResponse)
					break
				} else {
					// TODO: I am stuck here, how to send off-line msg?
					outMsg := OutMsg{TargetAddr:sendMsg.Target, Payload:sendMsg}
					log.Infof("[send] OutMsg send to msgQueue: %s, targetAddr: %s", outMsg, sendMsg.Target)
					msgQueue <- outMsg
				}
			case "hi":
				log.Infof("[hi] incoming>: %s", payload)
				// add this to clients pool
				hiMsg := models.HiMsg{}
				err := mapstructure.Decode(payload, &hiMsg)
				if err != nil {
					// send back to same address a error msg say, hi msg not right
					errorResponse := map[string]string{
						"error": "hi msg not right, should be have token, and user_addr field etc.",
					}
					ws.WriteJSON(errorResponse)
					break
				}else{
					// new clients online
					claims, err := utils.Decrypt(hiMsg.Token)
					if err != nil {
						// token invalid, expired or not right, refuse connection
						errorResponse := map[string]string{
							"error": "your token is invalid, re-login to refresh token.",
						}
						ws.WriteJSON(errorResponse)
						break
					} else{
						clientAddr := claims["user_addr"].(string)
						ua := hiMsg.UA
						client := Client{Conn: ws, ClientAddr:clientAddr, UA:ua, UpTime:time.Now().Unix()}
						// if clientAddr in clientsPool keys, then add to it's clients list, other wise not append
						if _, ok := clientsPool[clientAddr]; ok {
							clientsPool[clientAddr] = append(clientsPool[clientAddr], &client)
						} else {
							// clientAddr not in clientsPool
							clients := []*Client{&client}
							clientsPool[clientAddr] = clients
						}
						logOnlineDevicesInfo(clientsPool)
					}
				}
			case "add":
				log.Infof("[add] incoming>: %s", payload)
				continue
			case "del":
				log.Infof("[del] incoming>: %s", payload)
				continue
			default:
				errorResponse := map[string]string{
					"error": "only hi, send, add, del msg type are supported.",
				}
				ws.WriteJSON(errorResponse)
			}
		}
	}
}

func HandleMessages(){
	for {
		msg := <- msgQueue
		log.Info("msg from msgQueue: ", msg)
		if targetClients, ok := clientsPool[msg.TargetAddr]; ok {
			// one addr may have multi clients
			for i, client := range targetClients {
				err := client.Conn.WriteJSON(msg.Payload)
				if err != nil {
					log.Errorf("Error in broadcast msg. will drop this ws connection.")
					// drop lost connection
					clientsPool[msg.TargetAddr] = append(targetClients[:i], targetClients[i+1:]...)
				}
			}
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

func logOnlineDevicesInfo(clientsPool map[string][]*Client) {
	for k, v := range clientsPool {
		log.Infof("[Online] user_addr: %s  devices: %s", k, len(v))
	}
}

