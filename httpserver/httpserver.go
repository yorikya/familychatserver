package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/yorikya/familychatserver/client"
	"github.com/yorikya/familychatserver/hub"
)

//AuthResponse response to client who pass authentication
type AuthResponse struct {
	Resources string
}

func getURLParam(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("Url Param '%s' is missing", key)
	}
	return keys[0], nil
}

func broadcastHandler(h *hub.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Authenticate, add to room list from request ip
		msg, err := getURLParam(r, "msg")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}
		id, err := getURLParam(r, "id")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}

		h.BroadcastMessage(client.BroadcastMessage{
			Message: msg,
			UserID:  id,
		})
	}
}

func authHandler(h *hub.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Authenticate user
		id, err := getURLParam(r, "id")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}

		h.AddClient(&client.Client{
			ID: id,
			IP: strings.Split(r.RemoteAddr, ":")[0],
		})

		roomid, err := getURLParam(r, "roomid")
		if err != nil {
			log.Println(err)
			fmt.Fprintln(w, err)
			return
		}
		str, err := json.Marshal(AuthResponse{
			Resources: fmt.Sprintf("%s/rooms/%s/", h.GetResourcesPath(), roomid),
		})
		if err != nil {
			fmt.Fprintln(w, "failed decode response JSON"+err.Error())
		}
		fmt.Fprint(w, string(str))
	}
}

//Start starting the http and file system server
func Start(h *hub.Hub, port string) error {
	//File server handler
	directory := fmt.Sprintf(".%s", h.GetResourcesPath())
	log.Println("init static files server path:", directory)
	fs := http.FileServer(FileSystem{http.Dir(directory)})

	http.Handle(fmt.Sprintf("%s/", h.GetResourcesPath()), http.StripPrefix(h.GetResourcesPath(), fs))
	http.HandleFunc("/broadcast", broadcastHandler(h))
	http.HandleFunc("/auth", authHandler(h))

	log.Println("starting..., bind port:", port)
	return http.ListenAndServe(port, nil)
}
