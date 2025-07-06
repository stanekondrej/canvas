package canvas_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stanekondrej/canvas/server/internal/app/server/canvas"
)

func TestCanvasCreation(t *testing.T) {
	_ = canvas.NewCanvas()
}

func TestAddingStrokes(t *testing.T) {
	c := canvas.NewCanvas()

	s1 := canvas.Stroke{
		canvas.Coordinate{
			X: 1,
			Y: 8,
		},
	}
	c.AddStroke(s1)

	state := c.CurrentState()
	if len(state) != 1 {
		fmt.Println("Unexpected state length: expected 1, got", len(state))
		t.FailNow()
	}
}

func TestObservers(t *testing.T) {
	c := canvas.NewCanvas()

	receivedUpdate := false
	id := c.RegisterObserver(func(_s canvas.Stroke) {
		receivedUpdate = true
	})

	s := canvas.Stroke{
		canvas.Coordinate{
			X: 1,
			Y: 8,
		},
	}
	c.AddStroke(s)

	// allow the update to propagate
	time.Sleep(time.Millisecond * 50)

	if !receivedUpdate {
		fmt.Println("Update not received")
		t.FailNow()
	}

	c.UnregisterObserver(id)
}
