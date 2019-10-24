package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	id, ip string
}

type BroadcastMessage struct {
	Message, UserID string
}

func (c *Client) Send(m BroadcastMessage) error {
	url := fmt.Sprintf("http://%s:8080/message?id=%s&msg=%s", c.ip, m.UserID, url.QueryEscape(m.Message))
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("get response from client: %s, body: %s\n", c.ip, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Status code is %d", resp.StatusCode)
		return fmt.Errorf("Invalid Sttaus code: %d", resp.StatusCode)
	}

	fmt.Printf("Send msg to client: %s, from: %s. msg:%s\n", c.id, m.UserID, m.Message)
	return nil
}

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	fmt.Println("get request open file", path)
	return f, nil
}

type Hub struct {
	clients          map[string]*Client
	addClientChan    chan *Client
	removeClientChan chan *Client
	broadcastChan    chan BroadcastMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:          make(map[string]*Client),
		addClientChan:    make(chan *Client),
		removeClientChan: make(chan *Client),
		broadcastChan:    make(chan BroadcastMessage),
	}
}
func (h *Hub) removeClient(c *Client) {
	delete(h.clients, c.id)
}

func (h *Hub) addClient(c *Client) {
	fmt.Printf("client id: %s was added", c.id)
	h.clients[c.id] = c
}

func (h *Hub) broadcastMessage(m BroadcastMessage) {
	for _, client := range h.clients {
		err := client.Send(m)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			return
		}
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.addClientChan:
			h.addClient(conn)
		case conn := <-h.removeClientChan:
			h.removeClient(conn)
		case m := <-h.broadcastChan:
			h.broadcastMessage(m)
		}
	}
}

func main() {
	port := ":8080"
	hub := NewHub()
	fmt.Println("start, bind port:", port)

	http.HandleFunc("/bc", func(w http.ResponseWriter, r *http.Request) {
		//Authenticate, add to room list
		keys, ok := r.URL.Query()["msg"]
		if !ok || len(keys[0]) < 1 {
			rerr := "Url Param 'msg' is missing"
			log.Println(rerr)
			fmt.Fprintln(w, rerr)
			return
		}
		msg := keys[0]

		keys, ok = r.URL.Query()["id"]
		if !ok || len(keys[0]) < 1 {
			rerr := "Url Param 'id' is missing"
			log.Println(rerr)
			fmt.Fprintln(w, rerr)
			return
		}
		id := keys[0]
		// Query()["key"] will return an array of items,
		// we only want the single item.
		m := BroadcastMessage{
			Message: msg,
			UserID:  id,
		}
		hub.broadcastMessage(m)

	})

	resources := "/resources"
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		//Authenticate, add to room list
		keys, ok := r.URL.Query()["id"]
		if !ok || len(keys[0]) < 1 {
			rerr := "Url Param 'id' is missing"
			log.Println(rerr)
			fmt.Fprintln(w, rerr)
			return
		}

		// Query()["key"] will return an array of items,
		// we only want the single item.
		id := keys[0]
		fmt.Println("the client ip ", strings.Split(r.RemoteAddr, ":")[0])
		c := &Client{
			id: id,
			ip: strings.Split(r.RemoteAddr, ":")[0],
		}
		hub.addClient(c)

		fmt.Printf("the client ip: %s\n", r.RemoteAddr)
		m := make(map[string]string)
		m["resources"] = "/resource/rooms/1/"
		str, err := json.Marshal(m)
		fmt.Println("the response to client", string(str))
		if err != nil {
			fmt.Fprintln(w, "failed decode response JSON"+err.Error())
		}
		fmt.Fprint(w, string(str))

	})

	directory := fmt.Sprintf(".%s", resources)
	fmt.Println("init static files server path:", directory)
	fs := http.FileServer(FileSystem{http.Dir(directory)})
	http.Handle(fmt.Sprintf("%s/", resources), http.StripPrefix(resources, fs))

	log.Fatal(http.ListenAndServe(port, nil))
}
