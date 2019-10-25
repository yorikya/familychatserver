package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

//Client represent chat client
type Client struct {
	//ID client id
	ID,
	//IP client ip address
	IP,
	//Name client name
	Name string
}

// BroadcastMessage message from a client targeted to brodcasting to clients
type BroadcastMessage struct {
	//Message user message
	Message,
	//UserID user ID
	UserID string
}

//Send sends broadcast message to a client
func (c *Client) Send(m BroadcastMessage) error {
	url := fmt.Sprintf("http://%s:8080/message?id=%s&msg=%s", c.IP, m.UserID, url.QueryEscape(m.Message))
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("get response from client: %s, body: %s\n", c.IP, err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Status code is %d", resp.StatusCode)
		return fmt.Errorf("Invalid Sttaus code: %d", resp.StatusCode)
	}

	log.Printf("Send msg to client: %s, from: %s. msg:%s\n", c.ID, m.UserID, m.Message)
	return nil
}
