package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//MobileClient represent chat mobile client
type MobileClient struct {
	//ID client id
	ID,
	//IP client ip address
	IP,
	//Name client name
	Name string

	writeChan chan *BroadcastMessage
}

func (c *MobileClient) sendMessage(m *BroadcastMessage) error {
	url := fmt.Sprintf("http://%s:8080/message?id=%s&ts=%s&msg=%s", c.IP, m.UserID, m.Timestamp, url.QueryEscape(m.Message))
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

func NewMobileClient(id, ip, name string) *MobileClient {
	c := &MobileClient{
		ID:        id,
		IP:        ip,
		Name:      name,
		writeChan: make(chan *BroadcastMessage),
	}
	go c.run()

	return c
}

func (c *MobileClient) GetID() string {
	return c.ID
}

func (c *MobileClient) Close() {

}

//Send sends broadcast message to a client
func (c *MobileClient) Send(m *BroadcastMessage) {
	c.writeChan <- m
}
