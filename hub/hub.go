package hub

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/yorikya/familychatserver/client"
	"github.com/yorikya/familychatserver/db"
)

func getLastMessageIDFromLogFile(filePath string) int {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(fmt.Sprintf("%s%s", dir, filePath))
	if err != nil {
		log.Printf("failed to open log file '%s', error: %s", filePath, err)
		return 0
	}

	return len(strings.Split(string(content), "\n")) - 1
}

//Hub represent chat hub with clients
type Hub struct {
	clients          map[string]client.Client
	addClientChan    chan client.Client
	removeClientChan chan client.Client
	broadcastChan    chan *client.BroadcastMessage
	resourcePath     string
	dataBase         *db.DataBase
	roomID           string
	usersRoomPref    string
	lastMessageID    int
}

//NewHub return a new Hub
func NewHub(d *db.DataBase, roomID string) *Hub {
	h := &Hub{
		clients:          make(map[string]client.Client),
		addClientChan:    make(chan client.Client),
		removeClientChan: make(chan client.Client),
		broadcastChan:    make(chan *client.BroadcastMessage),
		resourcePath:     "/resources",
		dataBase:         d,
		roomID:           roomID,
		usersRoomPref:    "usersRoom",
		lastMessageID:    getLastMessageIDFromLogFile(fmt.Sprintf("/resources/rooms/%s/log.json", roomID)),
	}
	//Create users room bucket
	h.dataBase.CreateBucket(fmt.Sprintf("%s%s", h.usersRoomPref, h.roomID))

	go h.run()

	c, err := client.NewFileClient(roomID)
	if err != nil {
		panic(err)
	}
	h.AddClient(c)

	return h
}

func (h *Hub) Close() {
	for _, client := range h.clients {
		client.Close()
	}
}

func (h *Hub) GetLastMessageID() int {
	return h.lastMessageID
}

//AuthUser authenticated user password on success return true and error nil
func (h *Hub) AuthUser(user, pass string) (err error) {
	userPass, err := h.dataBase.GetAuthUser(user)
	if err != nil {
		return
	}
	if userPass != pass {
		err = fmt.Errorf("wrong password")
		return
	}
	return
}

//GetResourcesPath return resources path
func (h *Hub) GetResourcesPath() string {
	return h.resourcePath
}

func (h *Hub) removeClient(c client.Client) {
	log.Printf("client id: %s was deleted", c.GetID())
	delete(h.clients, c.GetID())
}

func (h *Hub) addClient(c client.Client) {
	log.Printf("client id: '%s' was added", c.GetID())
	h.clients[c.GetID()] = c
}

func (h *Hub) broadcastMessage(m *client.BroadcastMessage) {
	for _, client := range h.clients {
		go client.Send(m)
	}
}

//run starts comunication channel listeninig
func (h *Hub) run() {
	for {
		select {
		case c := <-h.addClientChan:
			h.addClient(c)
		case c := <-h.removeClientChan:
			h.removeClient(c)
		case m := <-h.broadcastChan:
			h.broadcastMessage(m)
		}
	}
}

//BroadcastMessage write message to the broadcast channel
func (h *Hub) BroadcastMessage(m *client.BroadcastMessage) {
	h.broadcastChan <- m
}

//AddClient adds client to the Hub
func (h *Hub) AddClient(c client.Client) {
	h.addClientChan <- c
}

//RemoveClient removes client from the Hub
func (h *Hub) RemoveClient(c client.Client) {
	h.removeClientChan <- c
}
