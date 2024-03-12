package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

var ug = websocket.Upgrader{
	HandshakeTimeout: 1 * time.Second,
	ReadBufferSize:   100,
	WriteBufferSize:  100,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub       *Hub
	conn      *websocket.Conn
	send      chan []byte
	frontName []byte
}

// reads from websocket conn, send to hub
func (client *Client) read() {
	defer func() { //graceful exit
		client.hub.unregister <- client //unregister hub from client
		log.Printf("%s offline \n", client.frontName)
		log.Printf("close connection to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close() //close websocket channel
	}()
	for {
		_, message, err := client.conn.ReadMessage() //when browser ends connection, err will be returned ，for loop breaks. then send chan closed in hub
		if err != nil {
			log.Printf("%v", err)
			break
		} else {
			message = bytes.TrimSpace(bytes.Replace(message, newLine, space, -1))
			if len(client.frontName) == 0 {
				client.frontName = message
				client.hub.broadcast <- []byte(fmt.Sprintf("%s online\n", string(client.frontName)))
			} else {
				client.hub.broadcast <- bytes.Join([][]byte{client.frontName, message}, []byte(": "))
			}
		}
	}
}

// read from broadcast ,write to websocket conn
func (client *Client) write() {
	defer func() {
		log.Printf("close connection to %s\n", client.conn.RemoteAddr().String())
		client.conn.Close() //close conn when fails to write
	}()

	for {
		msg, ok := <-client.send
		if !ok {
			log.Println("conn has been closed ")
			client.conn.WriteMessage(websocket.CloseMessage, []byte("bye bye"))
			return
		} else {
			err := client.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("failed to send msg to browser:%v\n", err)
				return
			}
		}
	}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := ug.Upgrade(w, r, nil) //http upgrade to websocket
	if err != nil {
		log.Printf("upgrade error: %v\n", err)
		return
	}
	log.Printf("connect to client %s\n", conn.RemoteAddr().String())
	// create a client for each web request
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	//register client to hub
	client.hub.register <- client

	//websocket is a full duplex protocol which reads and writes simultaneously.means it just works when concurrent.
	go client.read()
	go client.write()
}
