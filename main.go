package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}

type hub struct {
	clients          map[string]*websocket.Conn
	addClientChan    chan *websocket.Conn
	removeClientChan chan *websocket.Conn
	broadcastChan    chan Message
}

func (h *hub) removeClient(conn *websocket.Conn) {
	delete(h.clients, conn.LocalAddr().String())
}
func (h *hub) addClient(conn *websocket.Conn) {
	h.clients[conn.RemoteAddr().String()] = conn
}

func (h *hub) broadcastMessage(m Message) {
	for _, conn := range h.clients {
		err := websocket.JSON.Send(conn, m)
		if err != nil {
			fmt.Println("Error broadcasting message: ", err)
			return
		}
	}
}

func (h *hub) run() {
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

func newHub() *hub {
	return &hub{
		clients:          make(map[string]*websocket.Conn),
		addClientChan:    make(chan *websocket.Conn),
		removeClientChan: make(chan *websocket.Conn),
		broadcastChan:    make(chan Message),
	}
}

func handler(ws *websocket.Conn, h *hub) {
	go h.run()
	h.addClientChan <- ws
	for {
		var m Message
		err := websocket.JSON.Receive(ws, &m)
		if err != nil {
			h.broadcastChan <- Message{err.Error()}
			h.removeClient(ws)
			return
		}
		fmt.Printf("server side recive message %+v", m)
		h.broadcastChan <- m
	}
}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func server(port string) error {
	h := newHub()
	mux := http.NewServeMux()
	mux.Handle("/", websocket.Handler(func(ws *websocket.Conn) {
		handler(ws, h)
	}))

	mux.HandleFunc("/client", func(w http.ResponseWriter, r *http.Request) {
		//Authenticate..
		ip := mockedIP()
		fmt.Println("generate a new client:", ip)
		fmt.Fprintf(w, "%s", ip)
	})

	s := http.Server{Addr: ":" + port, Handler: mux}
	return s.ListenAndServe()
}

func main() {
	fmt.Println("hello")
	log.Fatal(server("8081"))
}
