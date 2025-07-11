package server

import (
	"encoding/json"
	"fmt"
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

	//go func() {
	//	for {
	//		if closed {
	//			break
	//		}

	//		cp := message.NewCheckpoint(h.canvas.CurrentState())
	//		if err := conn.WriteJSON(cp); err != nil {
	//			log.Println(err)
	//			continue
	//		}

	//		time.Sleep(CHECKPOINT_INTERVAL)
	//	}
	//}()

	// send initial canvas state
	{
		state := h.getCanvas().CurrentState()
		msg := message.NewCheckpoint(state)
		if err := conn.WriteJSON(msg); err != nil {

		}
	}

	// register interest in canvas updates
	id := h.getCanvas().RegisterObserver(func(s canvas.Stroke) {
		conn.WriteJSON(message.NewUpdate(s))
	})

	conn.SetCloseHandler(func(code int, text string) error {
		h.getCanvas().UnregisterObserver(id)

		if err := conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
			return err
		}

		return nil
	})

	for {
		t, p, err := conn.ReadMessage()
		if err != nil {
			conn.WriteJSON(message.NewError(err.Error()))
			continue
		}
		if t != websocket.TextMessage {
			conn.WriteJSON(message.NewError("Only text messages are supported"))
			continue
		}

		var msg struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(p, &msg); err != nil {
			conn.WriteJSON(message.NewError(fmt.Sprintf("Invalid JSON: %s", err)))
			continue
		}

		if msg.Type != "update" {
			conn.WriteJSON(message.NewError("Only updates are allowed"))
			continue
		}

		var data struct {
			Data canvas.Stroke `json:"data"`
		}
		if err := json.Unmarshal(p, &data); err != nil {
			conn.WriteJSON(message.NewError(fmt.Sprintf("Invalid stroke JSON: %s", err)))
			continue
		}

		h.getCanvas().AddStroke(data.Data)
	}
}
