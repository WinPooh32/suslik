package addon

import "github.com/WinPooh32/suslik"

type EventsBus struct {
	suslik.Game

	cur []interface{}
	new []interface{}
}

func NewEventsBus() *EventsBus {
	return &EventsBus{
		cur: make([]interface{}, 0, 128),
		new: make([]interface{}, 0, 128),
	}
}

func (bus *EventsBus) Send(event interface{}) {
	bus.new = append(bus.new, event)
}

func (bus *EventsBus) Events() []interface{} {
	return bus.cur
}

func (bus *EventsBus) EventsImmediate() []interface{} {
	return bus.new
}

func (bus *EventsBus) Reset() {
	bus.cur = bus.cur[:0]
	bus.new = bus.new[:0]
}

func (bus *EventsBus) Update(dt float32) {
	bus.next()
}

func (bus *EventsBus) next() {
	bus.cur, bus.new = bus.new, bus.cur[:0]
}
