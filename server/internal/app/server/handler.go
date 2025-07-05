package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stanekondrej/canvas/server/internal/app/server/canvas"
	"github.com/stanekondrej/canvas/server/internal/app/server/message"
)

const CHECKPOINT_INTERVAL time.Duration = time.Second * 10

type handler struct {
	canvas   canvas.Canvas
	upgrader websocket.Upgrader
}

func NewHandler() handler {
	return handler{
		canvas:   canvas.NewCanvas(),
		upgrader: getUpgrader(),
	}
}

func (h *handler) getCanvas() *canvas.Canvas {
	return &h.canvas
}

func getUpgrader() websocket.Upgrader {
	debugEnv, ok := os.LookupEnv("CANVAS_DEBUG")
	if !ok {
		return websocket.Upgrader{}
	}

	debug, err := strconv.ParseBool(debugEnv)
	if err != nil {
		return websocket.Upgrader{}
	}

	if !debug {
		return websocket.Upgrader{}
	}

	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (h *handler) Handle(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// send the current state of the canvas
	go func() {
		for {
			cp := message.NewCheckpoint(h.canvas.CurrentState())
			if err := conn.WriteJSON(cp); err != nil {
				log.Println(err)
				continue
			}

			time.Sleep(CHECKPOINT_INTERVAL)
		}
	}()

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

		if msg.MessageType == "close" {
			conn.WriteJSON(message.NewClose())
			break
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
