package message

import "github.com/stanekondrej/canvas/server/internal/app/server/canvas"

type errorMessage string

type Message struct {
	MessageType string `json:"type"`
	Data        any    `json:"data"`
}

func NewUpdate(s canvas.Stroke) Message {
	return Message{
		MessageType: "update",
		Data:        s,
	}
}

func NewCheckpoint(c []canvas.Stroke) Message {
	return Message{
		MessageType: "checkpoint",
		Data:        c,
	}
}

func NewError(e string) Message {
	return Message{
		MessageType: "error",
		Data:        errorMessage(e),
	}
}
