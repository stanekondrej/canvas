package canvas

import "sync"

type Coordinate struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Stroke []Coordinate

// a function able to handle updates
type ObserverFunc func(s Stroke)

type Canvas struct {
	// the current state of the canvas
	strokeLock sync.Mutex
	strokes    []Stroke

	observerID     ID
	observerIDLock sync.Mutex
	// functions waiting for an update
	observers map[ID]ObserverFunc
}

func NewCanvas() Canvas {
	return Canvas{
		strokes:   []Stroke{},
		observers: make(map[ID]ObserverFunc),
	}
}

func (c *Canvas) getObserverID() ID {
	var id ID

	c.observerIDLock.Lock()
	id = c.observerID
	c.observerID++
	c.observerIDLock.Unlock()

	return id
}

func (c *Canvas) CurrentState() []Stroke {
	return c.strokes
}

// unique identifier
type ID uint32

// adds a new observer to the canvas, for getting updates.
//
// returns an ID that can be later used to deregister the observer.
func (c *Canvas) RegisterObserver(o ObserverFunc) ID {
	id := c.getObserverID()
	c.observers[id] = o

	return id
}

func (c *Canvas) UnregisterObserver(id ID) {
	c.observers[id] = nil
}

// adds a stroke to the canvas
func (c *Canvas) AddStroke(s Stroke) {
	c.strokeLock.Lock()
	c.strokes = append(c.strokes, s)
	c.strokeLock.Unlock()

	go func() {
		for _, observer := range c.observers {
			go observer(s)
		}
	}()
}
