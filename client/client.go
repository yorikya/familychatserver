package client

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/yorikya/familychatserver/db"
)

//Client interface for chat clients
type Client interface {
	Send(m BroadcastMessage) error
	GetID() string
}

// BroadcastMessage message from a client targeted to brodcasting to clients
type BroadcastMessage struct {
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
}

func (c *MobileClient) GetID() string {
	return c.ID
}

//Send sends broadcast message to a client
func (c *MobileClient) Send(m BroadcastMessage) error {
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

//DBClient represent database client
type DBClient struct {
	ID     string
	msgNum int
	DB     *db.DataBase
}

func (c *DBClient) Send(m BroadcastMessage) error {
	err := c.DB.AddChatLog(fmt.Sprintf("%d", c.msgNum), fmt.Sprintf("%s|%s", m.UserID, m.Message))
	if err != nil {
		return err
	}
	c.msgNum++
	return nil
}

func (c *DBClient) GetID() string {
	return c.ID
}
