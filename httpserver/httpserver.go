package httpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/yorikya/familychatserver/client"
	"github.com/yorikya/familychatserver/hub"
)

//AuthResponse response to client who pass authentication
type authResponse struct {
	Error     string   `json:"error"`
	Resources string   `json:"resources"`
	Success   bool     `json:"success"`
	ClientID  string   `json:"clientid"`
	FileNames []string `json:"filenames"`
}

func newSuccessAuthResponse(resources, clientID string) authResponse {
	dir, err := os.Getwd()
	if err != nil {
		return newFailedAuthResponse(fmt.Errorf("failed to get pwd for client: '%s', error: %s", clientID, err))
	}

	files, err := ioutil.ReadDir(fmt.Sprintf("%s%s", dir, resources))
	if err != nil {
		return newFailedAuthResponse(fmt.Errorf("failed to get resources pass '%s', error: %s", resources, err))
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return authResponse{
		Resources: resources,
		Success:   true,
		ClientID:  clientID,
		FileNames: fileNames,
	}
}

func newFailedAuthResponse(err error) authResponse {
	return authResponse{
		Error: err.Error(),
	}
}

func handleAuthResponse(w http.ResponseWriter, auth authResponse) {
	str, err := json.Marshal(auth)
	if err != nil {
		fmt.Fprintln(w, fmt.Sprintf("failed decode response JSON, error: %s", err.Error()))
		return
	}
	log.Printf("response to client: %s", str)
	fmt.Fprint(w, string(str))
}
func getURLParam(r *http.Request, key string) (string, error) {
	keys, ok := r.URL.Query()[key]
	if !ok || len(keys[0]) < 1 {
		return "", fmt.Errorf("Url Param '%s' is missing", key)
	}
	return keys[0], nil
}

func broadcastHandler(h *hub.Hub) func(w http.ResponseWriter, r *http.Request) {
	msgID := h.GetLastMessageID()
	log.Printf("the last message id %d", msgID)
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

		h.BroadcastMessage(client.NewBroadcastMessage(msgID, id, msg))
		msgID++
	}
}

func authHandler(h *hub.Hub) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO: Authenticate user
		user, err := getURLParam(r, "user")
		if err != nil {
			handleAuthResponse(w, newFailedAuthResponse(err))
			return
		}

		pass, err := getURLParam(r, "pass")
		if err != nil {
			handleAuthResponse(w, newFailedAuthResponse(err))
			return
		}

		err = h.AuthUser(user, pass)
		if err != nil {
			handleAuthResponse(w, newFailedAuthResponse(err))
			return
		}

		h.AddClient(client.NewMobileClient(user, strings.Split(r.RemoteAddr, ":")[0], user))
		handleAuthResponse(w, newSuccessAuthResponse(fmt.Sprintf("%s/rooms/1/", h.GetResourcesPath()), user))
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
