package component

import (
	"github.com/WinPooh32/suslik"
)

type Bind struct {
	Action suslik.Action
	Keys   []suslik.Key
}

type BindAxis struct {
	Value float32
	Key   suslik.Key
}

type Input struct {
	Actions map[string]Bind
	Axes    map[string][]BindAxis

	keyboard    [512]byte
	mouseX      float32
	mouseY      float32
	mouseButton suslik.Key
	mouseAction suslik.Action
}

func MakeInput() Input {
	return Input{
		Actions: map[string]Bind{},
		Axes:    map[string][]BindAxis{},
	}
}

func (input *Input) MapAction(name string, bind Bind) {
	input.Actions[name] = bind
}

func (input *Input) MapAxis(name string, binds ...BindAxis) {
	input.Axes[name] = binds
}

func (input *Input) Action(name string) bool {
	b := input.Actions[name]
	for _, key := range b.Keys {
		if input.keyboard[key] != byte(b.Action) {
			return false
		}
	}
	return true
}

func (input *Input) Axis(name string) float32 {
	var sum float32
	var binds = input.Axes[name]

	for _, b := range binds {
		if input.keyboard[b.Key] != byte(suslik.NONE) {
			sum += b.Value
		}
	}
	return sum
}

func (input *Input) MouseState(name string) (x, y float32, button suslik.Key, action suslik.Action) {
	return input.mouseX, input.mouseY, input.mouseButton, input.mouseAction
}

func (input *Input) Update(dt float32) {
	for i := range input.keyboard {
		switch input.keyboard[i] {
		case byte(suslik.PRESS):
			input.keyboard[i] = byte(suslik.REPEAT)
		case byte(suslik.RELEASE):
			input.keyboard[i] = byte(suslik.NONE)
		}
	}
	input.mouseAction = suslik.NONE
}

func (input *Input) Mouse(x, y float32, button suslik.Key, action suslik.Action) {
	input.mouseX = x
	input.mouseY = y
	input.mouseButton = button
	input.mouseAction = action
}

func (input *Input) Scroll(amount float32) {

}

func (input *Input) Key(key suslik.Key, modifier suslik.Modifier, action suslik.Action) {
	input.keyboard[key] = byte(action)
}
