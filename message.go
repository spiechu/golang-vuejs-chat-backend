package main

import (
    "encoding/json"
    "fmt"
)

type OutgoingMessage struct {
    Type     Type   `json:"type"`
    Body     string `json:"body"`
    Sender   string `json:"sender"`
    SenderID int    `json:"senderID"`
}

type Type int

const (
    BackgroundUserInfo Type = iota
    UserMessage
    BackgroundSystemInfo
    UsersList
    UserInfo
)

func (m *OutgoingMessage) ToJSON() (stringified string, err error) {
    b, err := json.Marshal(m)

    if err != nil {
        fmt.Println(err)

        return "", err
    }

    fmt.Println(string(b))

    return string(b), nil
}
