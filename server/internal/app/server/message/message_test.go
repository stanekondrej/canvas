package message_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stanekondrej/canvas/server/internal/app/server/canvas"
	"github.com/stanekondrej/canvas/server/internal/app/server/message"
)

func compactJson(j []byte) (string, error) {
	var buf bytes.Buffer

	if err := json.Compact(&buf, j); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func TestCreateUpdate(t *testing.T) {
	stroke := canvas.Stroke{
		canvas.Coordinate{X: 1, Y: 4},
		canvas.Coordinate{X: 4, Y: 9},
	}
	msg := message.NewUpdate(stroke)

	j, err := json.Marshal(msg)
	if err != nil {
		t.FailNow()
	}

	e :=
		`{
			"type": "update",
			"data": [
				{ "x": 1, "y": 4 },
				{ "x": 4, "y": 9 }
			]
		}`
	expected, err := compactJson([]byte(e))
	if err != nil {
		fmt.Println("Unable to compact json:", err)

		t.FailNow()
	}

	if expected != string(j) {
		fmt.Println("Expected:\t", expected)
		fmt.Println("Got:\t", string(j))

		t.FailNow()
	}
}

func TestCreateError(t *testing.T) {
	msg := message.NewError("some error")

	j, err := json.Marshal(msg)
	if err != nil {
		t.FailNow()
	}

	e :=
		`{
			"type": "error",
			"data": "some error"
		}`

	expected, err := compactJson([]byte(e))
	if err != nil {
		fmt.Println("Unable to compact json:", err)

		t.FailNow()
	}

	if expected != string(j) {
		fmt.Println("Expected:\t", expected)
		fmt.Println("Got:\t", string(j))

		t.FailNow()
	}
}

func TestCreateCheckpoint(t *testing.T) {
	c := canvas.NewCanvas()
	c.AddStroke(canvas.Stroke{canvas.Coordinate{X: 1, Y: 8}})

	msg := message.NewCheckpoint(c.CurrentState())

	j, err := json.Marshal(msg)
	if err != nil {
		t.FailNow()
	}

	e :=
		`{
			"type": "checkpoint",
			"data": [
				[ 
					{
						"x": 1, 
						"y": 8 
					}
				]
			]
		}`

	expected, err := compactJson([]byte(e))
	if err != nil {
		fmt.Println("Unable to compact json:", err)

		t.FailNow()
	}

	if expected != string(j) {
		fmt.Println("Expected:\t", expected)
		fmt.Println("Got:\t", string(j))

		t.FailNow()
	}
}

func TestCreateCheckpointNoData(t *testing.T) {
	c := canvas.NewCanvas()

	msg := message.NewCheckpoint(c.CurrentState())

	j, err := json.Marshal(msg)
	if err != nil {
		t.FailNow()
	}

	e :=
		`{
			"type": "checkpoint",
			"data": []
		}`

	expected, err := compactJson([]byte(e))
	if err != nil {
		fmt.Println("Unable to compact json:", err)

		t.FailNow()
	}

	if expected != string(j) {
		fmt.Println("Expected:\t", expected)
		fmt.Println("Got:\t", string(j))

		t.FailNow()
	}
}

func TestCreateClose(t *testing.T) {
	msg := message.NewClose()

	j, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)

		t.FailNow()
	}

	e := `
	{
		"type": "close",
		"data": null
	}
	`

	expected, err := compactJson([]byte(e))
	if err != nil {
		fmt.Println(err)

		t.FailNow()
	}

	if expected != string(j) {
		fmt.Println("Expected:\t", expected)
		fmt.Println("Got:\t", string(j))

		t.FailNow()
	}
}
