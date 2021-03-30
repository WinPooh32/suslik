package suslik

type Responder interface {
	Render()
	Resize(width, height float32)
	Preload()
	Setup()
	Close()
	Update(dt float32)
	Mouse(x, y float32, button Mouse, action Action)
	Scroll(amount float32)
	Key(key Key, modifier Modifier, action Action)
	Type(char rune)
}

type Game struct{}

func (g *Game) Preload()                                        {}
func (g *Game) Setup()                                          {}
func (g *Game) Close()                                          {}
func (g *Game) Update(dt float32)                               {}
func (g *Game) Render()                                         {}
func (g *Game) Resize(w, h float32)                             {}
func (g *Game) Mouse(x, y float32, button Mouse, action Action) {}
func (g *Game) Scroll(amount float32)                           {}
func (g *Game) Key(key Key, modifier Modifier, action Action) {
	if key == Escape {
		Exit()
	}
}
func (g *Game) Type(char rune) {}
