package component

import "github.com/WinPooh32/suslik"

var (
	keyboardState [256]byte
	mouseX        float32
	mouseY        float32
	mouseButton   suslik.Key
	mouseAction   suslik.Action
)

type Bind struct {
	Action suslik.Action
	Keys   []suslik.Key
}

type Input struct {
	Actions map[string]Bind
}

func MakeInput() Input {
	return Input{
		Actions: map[string]Bind{},
	}
}

func (input *Input) MapAction(name string, action suslik.Action, keys ...suslik.Key) {
	input.Actions[name] = Bind{action, keys}
}

func (input *Input) MapAxis(name string) {
	// TODO
	// return 0
}

func (input *Input) Action(name string) bool {
	b := input.Actions[name]
	for _, key := range b.Keys {
		if keyboardState[key] != byte(b.Action) {
			return false
		}
	}
	return true
}

func (input *Input) Axis(name string) float32 {
	// TODO
	return 0
}

func (input *Input) MouseState(name string) (x, y float32, button suslik.Key, action suslik.Action) {
	return mouseX, mouseY, mouseButton, mouseAction
}

func (input *Input) Update(dt float32) {
	for i := range keyboardState {
		keyboardState[i] = byte(suslik.NONE)
	}
	mouseAction = suslik.NONE
}

func (input *Input) Mouse(x, y float32, button suslik.Key, action suslik.Action) {
	mouseX = x
	mouseY = y
	mouseButton = button
	mouseAction = action
}

func (input *Input) Scroll(amount float32) {

}

func (input *Input) Key(key suslik.Key, modifier suslik.Modifier, action suslik.Action) {
	keyboardState[key] = byte(action)
}
