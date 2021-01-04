package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

func getWsHandler(pool *Pool) websocket.Handler {
	return func(ws *websocket.Conn) {
		var clientID int

	OUTER:
		for {
			clientID = rand.Intn(1000000)

			if clientID == 0 {
				continue OUTER
			}

			for c := range pool.Clients {
				if clientID == c.ID {
					continue OUTER
				}
			}

			break OUTER
		}

		client := &Client{
			ID:   clientID,
			Name: fmt.Sprintf("Anonymous%d", clientID),
			Conn: ws,
			Pool: pool,
			RW:   &sync.RWMutex{},
		}

		pool.Register <- client

		client.Read()
	}
}

func main() {
	serverStartedAt := time.Now()

	pool := NewPool()

	go pool.Start()

	schedule := Schedule{
		Pool: pool,
	}

	schedule.Start(5*time.Second, "System health", func(s *Schedule) {
		ms := PrintMemUsage()

		ms.ServerStartedAt = serverStartedAt.Format(time.RFC822)
		ms.ActiveUsersCount = len(s.Pool.Clients)

		stringified, err := ms.ToJSON()
		if err != nil {
			fmt.Println(err)
			return
		}

		s.Pool.Broadcast <- OutgoingMessage{
			Type:   BackgroundSystemInfo,
			Body:   stringified,
			Sender: "",
		}
	})

	http.Handle("/", getWsHandler(pool))

	port := os.Getenv("PORT")

	if port == "" {
		port = "1234"
	}

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
