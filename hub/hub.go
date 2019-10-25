package hub

import (
	"fmt"
	"log"

	"github.com/yorikya/familychatserver/client"
	"github.com/yorikya/familychatserver/db"
)

//Hub represent chat hub with clients
type Hub struct {
	clients          map[string]client.Client
	addClientChan    chan client.Client
	removeClientChan chan client.Client
	broadcastChan    chan client.BroadcastMessage
	resourcePath     string
	dataBase         *db.DataBase
}

//NewHub return a new Hub
func NewHub(d *db.DataBase) *Hub {
	h := &Hub{
		clients:          make(map[string]client.Client),
		addClientChan:    make(chan client.Client),
		removeClientChan: make(chan client.Client),
		broadcastChan:    make(chan client.BroadcastMessage),
		resourcePath:     "/resources",
		dataBase:         d,
	}
	go h.run()

	h.AddClient(&client.DBClient{
		ID: "logfilewriter",
		DB: d,
	})
	return h
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
	log.Printf("client id: %s was added", c.GetID())
	h.clients[c.GetID()] = c
}

func (h *Hub) broadcastMessage(m client.BroadcastMessage) {
	for _, client := range h.clients {
		err := client.Send(m)
		if err != nil {
			log.Println("Error broadcasting message: ", err)
			return
		}
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
func (h *Hub) BroadcastMessage(m client.BroadcastMessage) {
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
