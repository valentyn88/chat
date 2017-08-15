package main

import (
	"bytes"
	"github.com/gorilla/websocket"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
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
		mtype, r, err := ws.NextReader()
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		if mtype == websocket.TextMessage {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r)
			broadcast <- buf.Bytes()
		} else if mtype == websocket.BinaryMessage {
			f, err := os.Create("../tmp/" + randStringRunes(10))
			if err != nil {
				log.Printf("error: %v", err)
				delete(clients, ws)
				break
			}
			defer f.Close()

			buf := make([]byte, 4096)
			for {
				n, err := r.Read(buf)
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Printf("error: %v", err)
					break
				}

				if n == 0 {
					break
				}

				if _, err := f.Write(buf[:n]); err != nil {
					log.Printf("error: %v", err)
					break
				}
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
