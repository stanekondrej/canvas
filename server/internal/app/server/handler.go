package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/stanekondrej/canvas/server/internal/app/server/canvas"
	"github.com/stanekondrej/canvas/server/internal/app/server/message"
)

type handler struct {
	canvas canvas.Canvas
}

func NewHandler() handler {
	return handler{
		canvas: canvas.NewCanvas(),
	}
}

func (h *handler) getCanvas() *canvas.Canvas {
	return &h.canvas
}

var upgrader websocket.Upgrader

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// send the current state of the canvas
	cp := message.NewCheckpoint(h.canvas.CurrentState())
	if err := conn.WriteJSON(cp); err != nil {
		log.Println(err)
		return
	}

	// register interest in canvas updates
	h.getCanvas().RegisterObserver(func(s canvas.Stroke) {
		conn.WriteJSON(message.NewUpdate(s))
	})

	for {
		var msg message.Message
		if err := conn.ReadJSON(&msg); err != nil {
			conn.WriteJSON(message.NewError(err.Error()))
			continue
		}
		if msg.MessageType != "update" {
			conn.WriteJSON(message.NewError("Only updates are allowed"))
			continue
		}

		// FIXME: this is so, so stupid
		j, _ := json.Marshal(msg)
		var incoming struct {
			Data canvas.Stroke `json:"data"`
		}
		if err := json.Unmarshal(j, &incoming); err != nil {
			conn.WriteJSON(message.NewError(err.Error()))
			continue
		}

		h.getCanvas().AddStroke(incoming.Data)
	}
}
