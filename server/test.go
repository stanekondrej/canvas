package main

import (
	"encoding/json"

	"github.com/stanekondrej/canvas/server/internal/app/server/canvas"
)

func main() {
	j, _ := json.Marshal(canvas.Stroke{
		canvas.Coordinate{X: 8, Y: 9},
	})

	print(string(j))
}
