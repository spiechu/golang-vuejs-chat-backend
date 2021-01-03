package main

import (
    "encoding/json"
    "fmt"
    "golang.org/x/net/websocket"
)

type SingleClientBroadcast struct {
    Client  *Client
    Message OutgoingMessage
}

type Pool struct {
    Register              chan *Client
    Unregister            chan *Client
    Clients               map[*Client]bool
    Broadcast             chan OutgoingMessage
    SingleClientBroadcast chan SingleClientBroadcast
}

func NewPool() *Pool {
    return &Pool{
        Register:              make(chan *Client),
        Unregister:            make(chan *Client),
        Clients:               make(map[*Client]bool),
        Broadcast:             make(chan OutgoingMessage),
        SingleClientBroadcast: make(chan SingleClientBroadcast),
    }
}

func (pool *Pool) Start() {
    for {
        select {
        case client := <-pool.Register:
            pool.Clients[client] = true

            fmt.Println("Registering")

            client.RW.RLock()
            message := OutgoingMessage{
                Type:     BackgroundUserInfo,
                Body:     "New User Joined...",
                Sender:   client.Name,
                SenderID: client.ID,
            }
            client.RW.RUnlock()

            userInfoMessage := OutgoingMessage{
                Type:     UserInfo,
                Body:     "Registered user",
                Sender:   client.Name,
                SenderID: client.ID,
            }

            sendMessageToOne(client, userInfoMessage)

            usersInfo, err := CreateUsersInfoMessage(pool)
            if err == nil {
                sendMessageToAll(pool, usersInfo)
            }

            sendMessageToAll(pool, message)

            break
        case client := <-pool.Unregister:
            delete(pool.Clients, client)

            fmt.Println("Unregistering")

            client.RW.RLock()
            message := OutgoingMessage{
                Type:     BackgroundUserInfo,
                Body:     "User Disconnected...",
                Sender:   client.Name,
                SenderID: client.ID,
            }
            client.RW.RUnlock()

            sendMessageToAll(pool, message)

            usersInfo, err := CreateUsersInfoMessage(pool)
            if err == nil {
                sendMessageToAll(pool, usersInfo)
            }

            break
        case message := <-pool.Broadcast:
            fmt.Println("Sending message to all clients in Pool")

            sendMessageToAll(pool, message)

            break
        case scb := <-pool.SingleClientBroadcast:
            fmt.Println("Sending message to single client")

            if _, ok := pool.Clients[scb.Client]; ok {
                sendMessageToOne(scb.Client, scb.Message)
            }

            break
        }
    }
}

func sendMessageToAll(pool *Pool, message OutgoingMessage) {
    fmt.Println("Size of Connection Pool: ", len(pool.Clients))

    stringified, err := message.ToJSON()
    if err != nil {
        fmt.Println("Cant stringify JSON, aborting broadcast")

        return
    }

    for client := range pool.Clients {
        if err = websocket.Message.Send(client.Conn, stringified); err != nil {
            fmt.Println("Can't send")
            continue
        }
    }
}

func sendMessageToOne(client *Client, message OutgoingMessage) {
    stringified, err := message.ToJSON()
    if err != nil {
        fmt.Println("Cant stringify JSON, aborting broadcast")

        return
    }

    if err = websocket.Message.Send(client.Conn, stringified); err != nil {
        fmt.Println("Can't send")
    }
}

func CreateUsersInfoMessage(pool *Pool) (OutgoingMessage, error) {
    var clients = map[int]string{}

    for client := range pool.Clients {
        clients[client.ID] = client.Name
    }

    b, err := json.Marshal(clients)

    if err != nil {
        fmt.Println(err)

        return OutgoingMessage{}, err
    }

    fmt.Println(string(b))

    return OutgoingMessage{Type: UsersList, Body: string(b), Sender: "", SenderID: 0}, nil
}
