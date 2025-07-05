package canvas

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Stroke []Coordinate

// a function able to handle updates
type ObserverFunc func(s Stroke)

type Canvas struct {
	// the current state of the canvas
	strokes []Stroke

	// functinos waiting for an update
	observers []ObserverFunc
}

func NewCanvas() Canvas {
	return Canvas{
		strokes:   []Stroke{},
		observers: []ObserverFunc{},
	}
}

func (c *Canvas) CurrentState() []Stroke {
	return c.strokes
}

// adds a new observer to the canvas, for getting updates.
func (c *Canvas) RegisterObserver(o ObserverFunc) {
	c.observers = append(c.observers, o)
}

// adds a stroke to the canvas
func (c *Canvas) AddStroke(s Stroke) {
	// TODO: consider moving this inside the goroutine
	c.strokes = append(c.strokes, s)

	go func() {
		for _, observer := range c.observers {
			go func() { observer(s) }()
		}
	}()
}
