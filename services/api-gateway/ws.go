package main

import (
	"log"
	"net/http"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/util"

	"github.com/gorilla/websocket"
)

var (
	userIDVar   = "userID"
	packageSlug = "packageSlug"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleRidersWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get(userIDVar)
	if len(userID) == 0 {
		log.Println("userID is not provided")
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read ws message failed %v", err)
			return
		}
		//conn.WriteJSON("ok")
		log.Printf("got ws message %s", string(msg))
	}
}

func handleDriverWebSocket(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed %v", err)
		return
	}

	defer conn.Close()

	userID := r.URL.Query().Get(userIDVar)
	if len(userID) == 0 {
		log.Println(userIDVar + "is not provide")
	}

	packageSlug := r.URL.Query().Get(packageSlug)
	if len(packageSlug) == 0 {
		log.Println(packageSlug + "is not provide")
	}

	//TODO refactor here
	type Driver struct {
		UserID         string `json:"userID"`
		Name           string `json:"name"`
		ProfilePicture string `json:"profilePicture"`
		CarPlate       string `json:"carPlate"`
		PackageSlug    string `json:"packageSlug"`
	}

	msg := contracts.WSMessage{
		Type: "driver.cmd.register",
		Data: Driver{
			UserID:         userID,
			Name:           "dodo",
			ProfilePicture: util.GetRandomAvatar(1),
			CarPlate:       "ABCDEFG123",
			PackageSlug:    packageSlug,
		},
	}

	if err := conn.WriteJSON(msg); err != nil {
		log.Printf("write ws response failed %v", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read ws message failed %v", err)
			return
		}
		log.Printf("got ws message %v", msg)
	}
}
