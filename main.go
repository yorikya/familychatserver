package main

import (
	"log"

	"github.com/yorikya/familychatserver/db"
	"github.com/yorikya/familychatserver/httpserver"
	"github.com/yorikya/familychatserver/hub"
)

func main() {
	port := ":8080"
	d := db.NewDB()
	h := hub.NewHub(d)
	defer h.Close()

	log.Fatal(httpserver.Start(h, port))
}
