package main

import (
	"encoding/json"

	"github.com/nats-io/nats.go"
)

func connectNats(url string) (*nats.Conn, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return nc, err
}

type BroadcastMsg struct {
	Id    int    `json:"id"`
	Msg   string `json:"msg"`
	Title string `json:"title"`
}

func NewBroadcastMsg(task Task, newState TaskState) *BroadcastMsg {
	return &BroadcastMsg{
		Id:    task.ID,
		Msg:   newState.String(),
		Title: task.Title,
	}
}

func (b *BroadcastMsg) Bytes() ([]byte, error) {
	return json.Marshal(b)
}
