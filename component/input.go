package component

import "github.com/WinPooh32/suslik"

var keyboardState [256]byte

type Bind struct {
	Action suslik.Action
	Keys   []suslik.Key
}

type Input struct {
	Actions map[string]Bind

	mouseX      float32
	mouseY      float32
	mouseButton suslik.Mouse
	mouseAction suslik.Action
}

func MakeInput() Input {
	return Input{
		Actions:     map[string]Bind{},
		mouseX:      0.0,
		mouseY:      0.0,
		mouseButton: 0,
		mouseAction: 0,
	}
}

func (input *Input) MapAction(name string, action suslik.Action, keys ...suslik.Key) {
	input.Actions[name] = Bind{action, keys}
}

func (input *Input) GetAction(name string) bool {
	b := input.Actions[name]
	for _, key := range b.Keys {
		if keyboardState[key] != byte(b.Action) {
			return false
		}
	}
	return true
}

func (input *Input) GetAxis(name string) float32 {
	// TODO
	return 0
}

func (input *Input) GetMouse(name string) (x, y float32, button suslik.Mouse, action suslik.Action) {
	return input.mouseX, input.mouseY, input.mouseButton, input.mouseAction
}

func (input *Input) Update(dt float32) {
	for i := range keyboardState {
		keyboardState[i] = byte(suslik.NONE)
	}
	input.mouseAction = suslik.NONE
}

func (input *Input) Mouse(x, y float32, button suslik.Mouse, action suslik.Action) {
	input.mouseX = x
	input.mouseY = y
	input.mouseButton = button
	input.mouseAction = action
}

func (input *Input) Scroll(amount float32) {

}

func (input *Input) Key(key suslik.Key, modifier suslik.Modifier, action suslik.Action) {
	keyboardState[key] = byte(action)
}
