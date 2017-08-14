## Gorilla websocket

go get github.com/gorilla/websocket

## Run server
cd server/
go run main.go

## Run client
cd client/
go run main.go

In client console start typing and press enter you will get message "Received: text message"

-a [path to file on local machine] - will store file with unique name to /tmp folder or get error message
if it's not a file or file larger than 200Mb
