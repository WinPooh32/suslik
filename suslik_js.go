// +build js

package suslik

import (
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/ajhager/webgl"
	"github.com/gopherjs/gopherjs/js"

	"github.com/WinPooh32/suslik/file"
)

func init() {
	rafPolyfill()
	rand.Seed(time.Now().UnixNano())
}

var canvas *js.Object
var keysDown map[int]bool

func run(title string, width, height int, fullscreen bool, hideCursor bool) {
	document := js.Global.Get("document")
	canvas = document.Call("createElement", "canvas")

	target := document.Call("getElementById", title)
	if target == nil {
		target = document.Get("body")
	}
	target.Call("appendChild", canvas)

	attrs := webgl.DefaultAttributes()
	attrs.Alpha = false
	attrs.Depth = false
	attrs.PremultipliedAlpha = false
	attrs.PreserveDrawingBuffer = false
	attrs.Antialias = false

	var err error
	gl, err = webgl.NewContext(canvas, attrs)
	if err != nil {
		log.Fatal(err)
	}

	js.Global.Set("onunload", func() {
		responder.Close()
	})

	canvas.Get("style").Set("display", "block")

	if hideCursor {
		canvas.Get("style").Set("cursor", "none")
	}

	canvas.Call("addEventListener", "contextmenu", func(ev *js.Object) {
		ev.Call("preventDefault")
		ev.Call("stopPropagation")
	}, false)

	winWidth := js.Global.Get("innerWidth").Int()
	winHeight := js.Global.Get("innerHeight").Int()
	if fullscreen {
		canvas.Set("width", winWidth)
		canvas.Set("height", winHeight)
	} else {
		canvas.Set("width", width)
		canvas.Set("height", height)
		canvas.Get("style").Set("marginLeft", toPx((winWidth-width)/2))
		canvas.Get("style").Set("marginTop", toPx((winHeight-height)/2))
	}

	canvas.Call("addEventListener", "mousemove", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		x := float32((ev.Get("clientX").Int() - rect.Get("left").Int()))
		y := float32((ev.Get("clientY").Int() - rect.Get("top").Int()))
		responder.Mouse(x, y, 0, MOVE)
	}, false)

	canvas.Call("addEventListener", "mousedown", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		x := float32((ev.Get("clientX").Int() - rect.Get("left").Int()))
		y := float32((ev.Get("clientY").Int() - rect.Get("top").Int()))
		btn := ev.Get("which").Int()
		responder.Mouse(x, y, toMouseBtn(btn), PRESS)
	}, false)

	canvas.Call("addEventListener", "mouseup", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		x := float32((ev.Get("clientX").Int() - rect.Get("left").Int()))
		y := float32((ev.Get("clientY").Int() - rect.Get("top").Int()))
		btn := ev.Get("which").Int()
		responder.Mouse(x, y, toMouseBtn(btn), RELEASE)
	}, false)

	canvas.Call("addEventListener", "touchstart", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, 0, PRESS)
		}
	}, false)

	canvas.Call("addEventListener", "touchcancel", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, 0, RELEASE)
		}
	}, false)

	canvas.Call("addEventListener", "touchend", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, 0, PRESS)
		}
	}, false)

	canvas.Call("addEventListener", "touchmove", func(ev *js.Object) {
		rect := canvas.Call("getBoundingClientRect")
		for i := 0; i < ev.Get("changedTouches").Get("length").Int(); i++ {
			touch := ev.Get("changedTouches").Index(i)
			x := float32((touch.Get("clientX").Int() - rect.Get("left").Int()))
			y := float32((touch.Get("clientY").Int() - rect.Get("top").Int()))
			responder.Mouse(x, y, 0, MOVE)
		}
	}, false)

	js.Global.Call("addEventListener", "keypress", func(ev *js.Object) {
		responder.Type(rune(ev.Get("charCode").Int()))
	}, false)

	keysDown = make(map[int]bool)

	js.Global.Call("addEventListener", "keydown", func(ev *js.Object) {
		key := ev.Get("keyCode").Int()
		if _, ok := keysDown[key]; ok {
			responder.Key(Key(key), 0, REPEAT)
		} else {
			keysDown[key] = true
			responder.Key(Key(key), 0, PRESS)
		}
	}, false)

	js.Global.Call("addEventListener", "keyup", func(ev *js.Object) {
		key := ev.Get("keyCode").Int()
		delete(keysDown, key)
		responder.Key(Key(key), 0, RELEASE)
	}, false)

	gl.Viewport(0, 0, width, height)

	responder.Preload()
	Files.Load(func() {
		responder.Setup()
		RequestAnimationFrame(animate)
	})
}

func width() float32 {
	return float32(canvas.Get("width").Int())
}

func height() float32 {
	return float32(canvas.Get("height").Int())
}

func animate(dt float32) {
	RequestAnimationFrame(animate)
	responder.Update(Time.Delta())
	gl.Clear(gl.COLOR_BUFFER_BIT)
	responder.Render()
	Time.Tick()
}

func exit() {
	responder.Close()
}

func toPx(n int) string {
	return strconv.FormatInt(int64(n), 10) + "px"
}

func toMouseBtn(b int) Key {
	switch b {
	case 1:
		return MouseLeft
	case 2:
		return MouseRight
	case 3:
		return MouseMiddle
	}
	return 0
}

func rafPolyfill() {
	window := js.Global
	vendors := []string{"ms", "moz", "webkit", "o"}
	if window.Get("requestAnimationFrame") == nil {
		for i := 0; i < len(vendors) && window.Get("requestAnimationFrame") == nil; i++ {
			vendor := vendors[i]
			window.Set("requestAnimationFrame", window.Get(vendor+"RequestAnimationFrame"))
			window.Set("cancelAnimationFrame", window.Get(vendor+"CancelAnimationFrame"))
			if window.Get("cancelAnimationFrame") == nil {
				window.Set("cancelAnimationFrame", window.Get(vendor+"CancelRequestAnimationFrame"))
			}
		}
	}

	lastTime := 0.0
	if window.Get("requestAnimationFrame") == nil {
		window.Set("requestAnimationFrame", func(callback func(float32)) int {
			currTime := js.Global.Get("Date").New().Call("getTime").Float()
			timeToCall := math.Max(0, 16-(currTime-lastTime))
			id := window.Call("setTimeout", func() { callback(float32(currTime + timeToCall)) }, timeToCall)
			lastTime = currTime + timeToCall
			return id.Int()
		})
	}

	if window.Get("cancelAnimationFrame") == nil {
		window.Set("cancelAnimationFrame", func(id int) {
			js.Global.Get("clearTimeout").Invoke(id)
		})
	}
}

func RequestAnimationFrame(callback func(float32)) int {
	return js.Global.Call("requestAnimationFrame", callback).Int()
}

func CancelAnimationFrame(id int) {
	js.Global.Call("cancelAnimationFrame")
}

func loadImage(r Resource) (Image, error) {
	ch := make(chan error, 1)

	img := js.Global.Get("Image").New()
	img.Call("addEventListener", "load", func(*js.Object) {
		go func() { ch <- nil }()
	}, false)
	img.Call("addEventListener", "error", func(o *js.Object) {
		go func() { ch <- &js.Error{Object: o} }()
	}, false)
	img.Set("src", r.url+"?"+strconv.FormatInt(rand.Int63(), 10))

	err := <-ch
	if err != nil {
		return nil, err
	}

	return &ImageObject{img}, nil
}

func loadJson(r Resource) (string, error) {
	ch := make(chan error, 1)

	req := js.Global.Get("XMLHttpRequest").New()
	req.Call("open", "GET", r.url, true)
	req.Call("addEventListener", "load", func(*js.Object) {
		go func() { ch <- nil }()
	}, false)
	req.Call("addEventListener", "error", func(o *js.Object) {
		go func() { ch <- &js.Error{Object: o} }()
	}, false)
	req.Call("send", nil)

	err := <-ch
	if err != nil {
		return "", err
	}

	return req.Get("responseText").String(), nil
}

func loadSound(r Resource) ([]byte, error) {
	return file.ReadAll(r.url)
}

type ImageObject struct {
	data *js.Object
}

func (i *ImageObject) Data() interface{} {
	return i.data
}

func (i *ImageObject) Width() int {
	return i.data.Get("width").Int()
}

func (i *ImageObject) Height() int {
	return i.data.Get("height").Int()
}
