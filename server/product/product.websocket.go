package product

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productSocket(ws *websocket.Conn) {
	done := make(chan int)
	go func(c *websocket.Conn) {
		fmt.Println("established a new websocket connection")
		for {
			var msg message
			if err := websocket.JSON.Receive(c, &msg); err != nil {
				log.Println(err)
				break
			}
			fmt.Printf("Recv %s", msg.Data)
		}
		close(done)
	}(ws)
loop:
	for {
		select {
		case <-done:
			fmt.Println("connection is closed. stopping.")
			break loop
		default:
			p, err := getTopTenProducts()
			if err != nil {
				log.Println(err)
				break
			}
			if err := websocket.JSON.Send(ws, p); err != nil {
				log.Println(err)
				break
			}
			time.Sleep(10 * time.Second)
		}
	}
	defer ws.Close()
}
