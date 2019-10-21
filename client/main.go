package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	Text string `json:"text"`
}

func mockedIP() string {
	var arr [4]int
	for i := 0; i < 4; i++ {
		rand.Seed(time.Now().UnixNano())
		arr[i] = rand.Intn(256)
	}
	return fmt.Sprintf("http://%d.%d.%d.%d", arr[0], arr[1], arr[2], arr[3])
}

func connect(ip string) (*websocket.Conn, error) {
	return websocket.Dial(fmt.Sprintf("ws://localhost:%s", "8081"), "", fmt.Sprintf("http://%s", ip))
}

func main() {
	resp, err := http.Get("http://localhost:8081/client")
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get IP from server", string(ip))

	ws, err := connect(string(ip))
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	// receive
	var m Message
	go func() {
		for {
			err := websocket.JSON.Receive(ws, &m)
			if err != nil {
				fmt.Println("Error receiving message: ", err.Error())

				connEstablish := false
				for !connEstablish {
					ws, err = connect(string(ip))
					if err != nil {
						log.Println("Error connect, sleep 1 sec", err)
						time.Sleep(1 * time.Second)
					} else {
						connEstablish = true
					}
				}

			}
			fmt.Println("Message: ", m)
		}
	}()
	// send
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}
		m := Message{
			Text: text,
		}
		err = websocket.JSON.Send(ws, m)
		if err != nil {
			fmt.Println("Error sending message: ", err.Error())
			break
		}
	}
}
