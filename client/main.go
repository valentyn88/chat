package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		url    = "ws://localhost:8000/ws"
		dialer *websocket.Dialer
		reader = bufio.NewReader(os.Stdin)
	)

	conn, _, err := dialer.Dial(url, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	go readMsg(conn)

	fmt.Println("Console chat started")

	for {
		//Read stdin
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		text = strings.Replace(text, "\n", "", -1)
		//Attach file
		if strings.Contains(text, "-a") {
			path := strings.Replace(text, "-a", "", -1)
			path = strings.TrimSpace(path)
			file, err := getFile(path)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			err = conn.WriteMessage(websocket.BinaryMessage, file)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
		} else if text != "" {
			err = conn.WriteMessage(websocket.TextMessage, []byte(text))
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
		}

	}
}

func readMsg(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}

		fmt.Printf("Received: %s\n", message)
	}
}

func getFile(path string) ([]byte, error) {
	var file []byte
	fInfo, err := os.Stat(path)
	if os.IsNotExist(err) || fInfo.IsDir() || fInfo.Size() > 200<<(10*2) {
		return nil, errors.New(fmt.Sprintf("%s is not a file or > 200Mb\n", path))
	} else {
		file, err = ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	return file, nil
}
