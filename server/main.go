package main

import (
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var (
	clients     = make(map[*websocket.Conn]bool)
	broadcast   = make(chan []byte)
	upgrader    = websocket.Upgrader{}
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()
	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	clients[ws] = true

	for {
		mtype, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		if mtype == websocket.TextMessage {
			broadcast <- message
		} else if mtype == websocket.BinaryMessage {
			err := ioutil.WriteFile("../tmp/"+randStringRunes(10), message, 0644)
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, ws)
				break
			}
		}
	}
}

func handleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
