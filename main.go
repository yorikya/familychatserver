package main

import (
	"log"

	"github.com/yorikya/familychatserver/httpserver"
	"github.com/yorikya/familychatserver/hub"
)

func main() {
	port := ":8080"
	h := hub.NewHub()
	go h.Run()

	log.Fatal(httpserver.Start(h, port))
}
