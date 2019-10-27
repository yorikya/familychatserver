package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type FileClient struct {
	ID        string
	file      *os.File
	roomID    string
	writeChan chan *BroadcastMessage
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

func (c FileClient) writeLine(m *BroadcastMessage) error {
	str, err := json.Marshal(m)
	if err != nil {
		return err
	}
	c.file.WriteString(fmt.Sprintf("%s\n", string(str)))
	return nil
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
		writeChan: make(chan *BroadcastMessage),
		roomID:    roomID,
		file:      f,
		ID:        fmt.Sprintf("log-file-writer-room-%s", roomID),
	}

	go c.run()

	return c, nil
}

func (c *FileClient) Send(m *BroadcastMessage) {
	c.writeChan <- m
}

func (c *FileClient) GetID() string {
	return c.ID
}

func (c *FileClient) Close() {
	c.file.Close()
}
