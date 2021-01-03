package main

import (
    "encoding/json"
    "fmt"
    "golang.org/x/net/websocket"
    "sync"
)

type Client struct {
    ID   int
    Name string
    Conn *websocket.Conn
    RW   *sync.RWMutex
    Pool *Pool
}

type IncomingMessage struct {
    Command string
    Body    string
}

func (c *Client) Read() {
    defer func() {
        fmt.Println("Closing Client Read")

        c.Pool.Unregister <- c
        _ = c.Conn.Close()
    }()

    var (
        err        error
        rawMessage string
    )

    for {
        err = websocket.Message.Receive(c.Conn, &rawMessage)
        if err != nil {
            fmt.Println("Can't receive")

            // fatal error, connection broken
            break
        }

        fmt.Println("Received from client: " + rawMessage)

        incomingMessage := &IncomingMessage{}

        err = json.Unmarshal([]byte(rawMessage), incomingMessage)
        if err != nil {
            fmt.Println(err)

            // non fatal error, not disconnecting client
            continue
        }

        switch incomingMessage.Command {
        case "setName":
            c.RW.Lock()
            oldName := c.Name
            c.Name = incomingMessage.Body

            userInfoMessage := OutgoingMessage{
                Type:     UserInfo,
                Body:     "Registered user",
                Sender:   c.Name,
                SenderID: c.ID,
            }

            message := OutgoingMessage{
                Type:     BackgroundUserInfo,
                Body:     "User " + oldName + " changed name to " + c.Name,
                Sender:   c.Name,
                SenderID: c.ID,
            }
            c.RW.Unlock()

            c.Pool.SingleClientBroadcast <- SingleClientBroadcast{
                Client:  c,
                Message: userInfoMessage,
            }

            c.Pool.Broadcast <- message

            usersInfo, err := CreateUsersInfoMessage(c.Pool)
            if err == nil {
                sendMessageToAll(c.Pool, usersInfo)
            }
        case "sendMessage":
            c.RW.RLock()
            message := OutgoingMessage{Type: UserMessage, Body: incomingMessage.Body, Sender: c.Name, SenderID: c.ID}
            c.RW.RUnlock()

            c.Pool.Broadcast <- message
        default:
            fmt.Println("Command unrecognized: " + incomingMessage.Command)
        }
    }
}
