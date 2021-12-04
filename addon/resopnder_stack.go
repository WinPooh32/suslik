package addon

import "github.com/WinPooh32/suslik"

type ResponderStack struct {
	stack []suslik.Responder
}

func NewResponderStack(responder ...suslik.Responder) *ResponderStack {
	if len(responder) == 0 {
		panic("responder list must not be empty")
	}
	return &ResponderStack{
		stack: responder,
	}
}

func (rs *ResponderStack) Render() {
	for _, r := range rs.stack {
		r.Render()
	}
}

func (rs *ResponderStack) Resize(width, height float32) {
	for _, r := range rs.stack {
		r.Resize(width, height)
	}
}

func (rs *ResponderStack) Preload() {
	for _, r := range rs.stack {
		r.Preload()
	}
}

func (rs *ResponderStack) Setup() {
	for _, r := range rs.stack {
		r.Setup()
	}
}

func (rs *ResponderStack) Close() {
	for _, r := range rs.stack {
		r.Close()
	}
}

func (rs *ResponderStack) Update(dt float32) {
	for _, r := range rs.stack {
		r.Update(dt)
	}
}

func (rs *ResponderStack) Mouse(x, y float32, button suslik.Key, action suslik.Action) {
	for _, r := range rs.stack {
		r.Mouse(x, y, button, action)
	}
}

func (rs *ResponderStack) Scroll(amount float32) {
	for _, r := range rs.stack {
		r.Scroll(amount)
	}
}

func (rs *ResponderStack) Key(key suslik.Key, modifier suslik.Modifier, action suslik.Action) {
	for _, r := range rs.stack {
		r.Key(key, modifier, action)
	}
}

func (rs *ResponderStack) Type(char rune) {
	for _, r := range rs.stack {
		r.Type(char)
	}
}
