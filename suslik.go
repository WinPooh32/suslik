package suslik

import "github.com/ajhager/webgl"

var (
	responder Responder
	Time      *Clock
	Files     *Loader
	gl        *webgl.Context
)

func Open(title string, width, height int, fullscreen bool, hideCursor bool, r Responder) {
	responder = r
	Time = NewClock()
	Files = NewLoader()
	run(title, width, height, fullscreen, hideCursor)
}

func SetBg(color uint32) {
	r := float32((color>>16)&0xFF) / 255.0
	g := float32((color>>8)&0xFF) / 255.0
	b := float32(color&0xFF) / 255.0
	gl.ClearColor(r, g, b, 1.0)
}

func Width() float32 {
	return width()
}

func Height() float32 {
	return height()
}

func Exit() {
	exit()
}
