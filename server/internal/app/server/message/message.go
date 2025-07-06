package message

import "github.com/stanekondrej/canvas/server/internal/app/server/canvas"

type errorMessage string

type MessageType string

const (
	UpdateMessage     MessageType = "update"
	CheckpointMessage MessageType = "checkpoint"
	ErrorMessage      MessageType = "error"
)

type Message[T errorMessage | []canvas.Stroke | canvas.Stroke] struct {
	Type MessageType `json:"type"`
	Data T           `json:"data"`
}

func NewUpdate(s canvas.Stroke) Message[canvas.Stroke] {
	return Message[canvas.Stroke]{
		Type: "update",
		Data: s,
	}
}

func NewCheckpoint(c []canvas.Stroke) Message[[]canvas.Stroke] {
	return Message[[]canvas.Stroke]{
		Type: "checkpoint",
		Data: c,
	}
}

func NewError(e string) Message[errorMessage] {
	return Message[errorMessage]{
		Type: "error",
		Data: errorMessage(e),
	}
}
