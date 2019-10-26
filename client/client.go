package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

//Client interface for chat clients
type Client interface {
	Send(m BroadcastMessage)
	GetID() string
	Close()
}

// BroadcastMessage message from a client targeted to brodcasting to clients
type BroadcastMessage struct {
	MessageID int
	//Message user message
	Message,
	//UserID user ID
	UserID string
}

//MobileClient represent chat mobile client
type MobileClient struct {
	//ID client id
	ID,
	//IP client ip address
	IP,
	//Name client name
	Name string

	writeChan chan BroadcastMessage
}

func NewMobileClient(id, ip, name string) *MobileClient {
	c := &MobileClient{
		ID:        id,
		IP:        ip,
		Name:      name,
		writeChan: make(chan BroadcastMessage),
	}
	go c.run()

	return c
}

func (c *MobileClient) sendMessage(m BroadcastMessage) error {
	url := fmt.Sprintf("http://%s:8080/message?id=%s&msg=%s", c.IP, m.UserID, url.QueryEscape(m.Message))
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("get error from client: %s, error: %s", c.IP, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("get error status code from client: %s, status: %d", c.IP, resp.StatusCode)
	}

	log.Printf("Send msg to client: %s, from: %s. msg:%s\n", c.ID, m.UserID, m.Message)
	return nil
}

func (c *MobileClient) run() {
	for m := range c.writeChan {
		err := c.sendMessage(m)
		if err != nil {
			log.Println("failed send messge to file error:", err)
		}
	}
}

func (c *MobileClient) GetID() string {
	return c.ID
}

func (c *MobileClient) Close() {

}

//Send sends broadcast message to a client
func (c *MobileClient) Send(m BroadcastMessage) {
	c.writeChan <- m
}

func NewFileClient(roomID string) (*FileClient, error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf("%s/resources/rooms/%s/log.json", dir, roomID)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return nil, err
	}
	log.Println("FileClient creates file :", path)
	c := &FileClient{
		writeChan: make(chan BroadcastMessage),
		roomID:    roomID,
		file:      f,
		ID:        fmt.Sprintf("log-file-writer-room-%s", roomID),
	}

	go c.run()

	return c, nil
}

type FileClient struct {
	ID        string
	msgNum    int
	file      *os.File
	roomID    string
	writeChan chan BroadcastMessage
}

//run starts comunication channel listeninig
func (c *FileClient) run() {
	for m := range c.writeChan {
		err := c.writeLine(m)
		if err != nil {
			log.Println("failed send messge to file error:", err)
		}
	}
}

func (c FileClient) writeLine(m BroadcastMessage) error {
	m.MessageID = c.msgNum
	str, err := json.Marshal(m)
	if err != nil {
		return err
	}
	c.file.WriteString(fmt.Sprintf("%s\n", string(str)))
	c.msgNum++
	return nil
}
func (c *FileClient) Send(m BroadcastMessage) {
	c.writeChan <- m
}

func (c *FileClient) GetID() string {
	return c.ID
}

func (c *FileClient) Close() {
	c.file.Close()
}
