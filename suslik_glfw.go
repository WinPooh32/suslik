// +build !js

package suslik

import (
	"image"
	"image/draw"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/ajhager/webgl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var window *glfw.Window

// fatalErr calls log.Fatal with the given error if it is non-nil.
func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	windowWidth  int
	windowHeight int
)

func init() {
	runtime.LockOSThread()
}

func run(title string, width, height int, fullscreen bool, hideCursor bool) {
	fatalErr(glfw.Init())

	monitor := glfw.GetPrimaryMonitor()
	mode := monitor.GetVideoMode()

	if fullscreen {
		width = mode.Width
		height = mode.Height
		glfw.WindowHint(glfw.Decorated, glfw.False)
	} else {
		monitor = nil
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(width, height, title, monitor, nil)
	fatalErr(err)
	window.MakeContextCurrent()

	if !fullscreen {
		window.SetPos((mode.Width-width)/2, (mode.Height-height)/2)
	}

	if hideCursor {
		window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	}

	windowWidth, windowHeight = window.GetSize()
	width, height = window.GetFramebufferSize()

	glfw.SwapInterval(1)

	gl = webgl.NewContext()

	gl.Viewport(0, 0, width, height)
	window.SetFramebufferSizeCallback(func(window *glfw.Window, w, h int) {
		windowWidth, windowHeight = window.GetSize()
		gl.Viewport(0, 0, w, h)
		responder.Resize(float32(windowWidth), float32(windowHeight))
	})

	window.SetCursorPosCallback(func(window *glfw.Window, x, y float64) {
		responder.Mouse(float32(x), float32(y), 0, MOVE)
	})

	window.SetMouseButtonCallback(func(window *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		x, y := window.GetCursorPos()

		var mb Key
		switch b {
		case glfw.MouseButtonLeft:
			mb = MouseLeft
		case glfw.MouseButtonMiddle:
			mb = MouseMiddle
		case glfw.MouseButtonRight:
			mb = MouseRight
		}

		if a == glfw.Press {
			responder.Mouse(float32(x), float32(y), mb, PRESS)
		} else {
			responder.Mouse(float32(x), float32(y), mb, RELEASE)
		}
	})

	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		responder.Scroll(float32(yoff))
	})

	window.SetKeyCallback(func(window *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		switch a {
		case glfw.Press:
			responder.Key(Key(k), Modifier(m), PRESS)
		case glfw.Release:
			responder.Key(Key(k), Modifier(m), RELEASE)
		case glfw.Repeat:
			responder.Key(Key(k), Modifier(m), REPEAT)
		}
	})

	window.SetCharCallback(func(window *glfw.Window, char rune) {
		responder.Type(char)
	})

	responder.Preload()
	Files.Load(func() {})
	responder.Setup()

	shouldClose := window.ShouldClose()
	for !shouldClose {
		responder.Update(Time.Delta())
		gl.Clear(gl.COLOR_BUFFER_BIT)
		responder.Render()
		window.SwapBuffers()
		glfw.PollEvents()
		Time.Tick()

		shouldClose = window.ShouldClose()
	}
	window.Destroy()
	glfw.Terminate()
	responder.Close()
}

func width() float32 {
	return float32(windowWidth)
}

func height() float32 {
	return float32(windowHeight)
}

func exit() {
	window.SetShouldClose(true)
}

func init() {
	BoardDash = Key(glfw.KeyMinus)
	BoardApostrophe = Key(glfw.KeyApostrophe)
	BoardSemicolon = Key(glfw.KeySemicolon)
	BoardEquals = Key(glfw.KeyEqual)
	BoardComma = Key(glfw.KeyComma)
	BoardPeriod = Key(glfw.KeyPeriod)
	BoardSlash = Key(glfw.KeySlash)
	BoardBackslash = Key(glfw.KeyBackslash)
	BoardBackspace = Key(glfw.KeyBackspace)
	BoardTab = Key(glfw.KeyTab)
	BoardCapsLock = Key(glfw.KeyCapsLock)
	BoardSpace = Key(glfw.KeySpace)
	BoardEnter = Key(glfw.KeyEnter)
	BoardEscape = Key(glfw.KeyEscape)
	BoardInsert = Key(glfw.KeyInsert)
	BoardPrintScreen = Key(glfw.KeyPrintScreen)
	BoardDelete = Key(glfw.KeyDelete)
	BoardPageUp = Key(glfw.KeyPageUp)
	BoardPageDown = Key(glfw.KeyPageDown)
	BoardHome = Key(glfw.KeyHome)
	BoardEnd = Key(glfw.KeyEnd)
	BoardPause = Key(glfw.KeyPause)
	BoardScrollLock = Key(glfw.KeyScrollLock)
	BoardArrowLeft = Key(glfw.KeyLeft)
	BoardArrowRight = Key(glfw.KeyRight)
	BoardArrowDown = Key(glfw.KeyDown)
	BoardArrowUp = Key(glfw.KeyUp)
	BoardLeftBracket = Key(glfw.KeyLeftBracket)
	BoardLeftShift = Key(glfw.KeyLeftShift)
	BoardLeftControl = Key(glfw.KeyLeftControl)
	BoardLeftSuper = Key(glfw.KeyLeftSuper)
	BoardLeftAlt = Key(glfw.KeyLeftAlt)
	BoardRightBracket = Key(glfw.KeyRightBracket)
	BoardRightShift = Key(glfw.KeyRightShift)
	BoardRightControl = Key(glfw.KeyRightControl)
	BoardRightSuper = Key(glfw.KeyRightSuper)
	BoardRightAlt = Key(glfw.KeyRightAlt)
	BoardZero = Key(glfw.Key0)
	BoardOne = Key(glfw.Key1)
	BoardTwo = Key(glfw.Key2)
	BoardThree = Key(glfw.Key3)
	BoardFour = Key(glfw.Key4)
	BoardFive = Key(glfw.Key5)
	BoardSix = Key(glfw.Key6)
	BoardSeven = Key(glfw.Key7)
	BoardEight = Key(glfw.Key8)
	BoardNine = Key(glfw.Key9)
	BoardF1 = Key(glfw.KeyF1)
	BoardF2 = Key(glfw.KeyF2)
	BoardF3 = Key(glfw.KeyF3)
	BoardF4 = Key(glfw.KeyF4)
	BoardF5 = Key(glfw.KeyF5)
	BoardF6 = Key(glfw.KeyF6)
	BoardF7 = Key(glfw.KeyF7)
	BoardF8 = Key(glfw.KeyF8)
	BoardF9 = Key(glfw.KeyF9)
	BoardF10 = Key(glfw.KeyF10)
	BoardF11 = Key(glfw.KeyF11)
	BoardF12 = Key(glfw.KeyF12)
	BoardA = Key(glfw.KeyA)
	BoardB = Key(glfw.KeyB)
	BoardC = Key(glfw.KeyC)
	BoardD = Key(glfw.KeyD)
	BoardE = Key(glfw.KeyE)
	BoardF = Key(glfw.KeyF)
	BoardG = Key(glfw.KeyG)
	BoardH = Key(glfw.KeyH)
	BoardI = Key(glfw.KeyI)
	BoardJ = Key(glfw.KeyJ)
	BoardK = Key(glfw.KeyK)
	BoardL = Key(glfw.KeyL)
	BoardM = Key(glfw.KeyM)
	BoardN = Key(glfw.KeyN)
	BoardO = Key(glfw.KeyO)
	BoardP = Key(glfw.KeyP)
	BoardQ = Key(glfw.KeyQ)
	BoardR = Key(glfw.KeyR)
	BoardS = Key(glfw.KeyS)
	BoardT = Key(glfw.KeyT)
	BoardU = Key(glfw.KeyU)
	BoardV = Key(glfw.KeyV)
	BoardW = Key(glfw.KeyW)
	BoardX = Key(glfw.KeyX)
	BoardY = Key(glfw.KeyY)
	BoardZ = Key(glfw.KeyZ)
	BoardNumLock = Key(glfw.KeyNumLock)
	BoardNumMultiply = Key(glfw.KeyKPMultiply)
	BoardNumDivide = Key(glfw.KeyKPDivide)
	BoardNumAdd = Key(glfw.KeyKPAdd)
	BoardNumSubtract = Key(glfw.KeyKPSubtract)
	BoardNumZero = Key(glfw.KeyKP0)
	BoardNumOne = Key(glfw.KeyKP1)
	BoardNumTwo = Key(glfw.KeyKP2)
	BoardNumThree = Key(glfw.KeyKP3)
	BoardNumFour = Key(glfw.KeyKP4)
	BoardNumFive = Key(glfw.KeyKP5)
	BoardNumSix = Key(glfw.KeyKP6)
	BoardNumSeven = Key(glfw.KeyKP7)
	BoardNumEight = Key(glfw.KeyKP8)
	BoardNumNine = Key(glfw.KeyKP9)
	BoardNumDecimal = Key(glfw.KeyKPDecimal)
	BoardNumEnter = Key(glfw.KeyKPEnter)
}

func NewImageObject(img *image.NRGBA) *ImageObject {
	return &ImageObject{img}
}

type ImageObject struct {
	data *image.NRGBA
}

func (i *ImageObject) Data() interface{} {
	return i.data
}

func (i *ImageObject) Width() int {
	return i.data.Rect.Max.X
}

func (i *ImageObject) Height() int {
	return i.data.Rect.Max.Y
}

func loadImage(r Resource) (Image, error) {
	file, err := os.Open(r.url)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	b := img.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), img, b.Min, draw.Src)

	return &ImageObject{newm}, nil
}

func loadJson(r Resource) (string, error) {
	file, err := ioutil.ReadFile(r.url)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func loadSound(r Resource) ([]byte, error) {
	return ioutil.ReadFile(r.url)
}

type Assets struct {
	queue  []string
	cache  map[string]Image
	loads  int
	errors int
}

func NewAssets() *Assets {
	return &Assets{make([]string, 0), make(map[string]Image), 0, 0}
}

func (a *Assets) Image(path string) {
	a.queue = append(a.queue, path)
}

func (a *Assets) Get(path string) Image {
	return a.cache[path]
}

func (a *Assets) Load(onFinish func()) {
	if len(a.queue) == 0 {
		onFinish()
	} else {
		for _, path := range a.queue {
			img := LoadImage(path)
			a.cache[path] = img
		}
	}
}

func LoadImage(data interface{}) Image {
	var m image.Image

	switch data := data.(type) {
	default:
		log.Fatal("NewTexture needs a string or io.Reader")
	case string:
		file, err := os.Open(data)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		img, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case io.Reader:
		img, _, err := image.Decode(data)
		if err != nil {
			log.Fatal(err)
		}
		m = img
	case image.Image:
		m = data
	}

	b := m.Bounds()
	newm := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), m, b.Min, draw.Src)

	return &ImageObject{newm}
}
