package suslik

import (
	"log"
	"path"
	"strings"

	"github.com/WinPooh32/suslik/audio"
	"github.com/WinPooh32/suslik/math"
	"github.com/ajhager/webgl"
)

type Resource struct {
	kind string
	name string
	url  string
}

type Loader struct {
	resources []Resource
	images    map[string]*Texture
	jsons     map[string]string
	sounds    map[string]*Sound
}

func NewLoader() *Loader {
	return &Loader{
		resources: make([]Resource, 1),
		images:    make(map[string]*Texture),
		jsons:     make(map[string]string),
		sounds:    make(map[string]*Sound),
	}
}

func (l *Loader) Add(name, url string) {
	kind := strings.ToLower(path.Ext(url))
	l.resources = append(l.resources, Resource{kind, name, url})
}

func (l *Loader) Image(name string) *Texture {
	return l.images[name]
}

func (l *Loader) Json(name string) string {
	return l.jsons[name]
}

func (l *Loader) Sound(name string) *Sound {
	return l.sounds[name]
}

func (l *Loader) Load(onFinish func()) {
	for _, r := range l.resources {
		switch r.kind {
		case ".png":
			data, err := loadImage(r)
			if err != nil {
				log.Println(r.url, "png load:", err)
				continue
			}

			l.images[r.name] = NewTexture(data)

		case ".json":
			data, err := loadJson(r)
			if err != nil {
				log.Println(r.url, "json load:", err)
				continue
			}

			l.jsons[r.name] = data

		case ".wav", ".mp3", ".ogg":
			data, err := loadSound(r)
			if err != nil {
				log.Printf("load sound: %s: %s", r.url, err)
				continue
			}

			decoded, err := audio.Decode(data, r.kind)
			if err != nil {
				log.Printf("decode sound: %s: %s", r.url, err)
				continue
			}

			l.sounds[r.name] = &Sound{audio.NewPlayer(decoded)}
		}
	}
	onFinish()
}

type Image interface {
	Data() interface{}
	Width() int
	Height() int
}

func LoadShader(vertSrc, fragSrc string) *webgl.Program {
	vertShader := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertShader, vertSrc)
	gl.CompileShader(vertShader)
	defer gl.DeleteShader(vertShader)

	fragShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragShader, fragSrc)
	gl.CompileShader(fragShader)
	defer gl.DeleteShader(fragShader)

	program := gl.CreateProgram()
	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)
	gl.LinkProgram(program)

	return program
}

type Region struct {
	texture       *Texture
	u, v          float32
	u2, v2        float32
	width, height float32
}

func NewRegion(texture *Texture, x, y, w, h int) *Region {
	invTexWidth := 1.0 / float32(texture.Width())
	invTexHeight := 1.0 / float32(texture.Height())

	u := float32(x) * invTexWidth
	v := float32(y) * invTexHeight
	u2 := float32(x+w) * invTexWidth
	v2 := float32(y+h) * invTexHeight
	width := float32(math.Abs(float32(w)))
	height := float32(math.Abs(float32(h)))

	return &Region{texture, u, v, u2, v2, width, height}
}

func (r *Region) Width() float32 {
	return float32(r.width)
}

func (r *Region) Height() float32 {
	return float32(r.height)
}

func (r *Region) Texture() *webgl.Texture {
	return r.texture.id
}

func (r *Region) View() (float32, float32, float32, float32) {
	return r.u, r.v, r.u2, r.v2
}

type Texture struct {
	id     *webgl.Texture
	width  int
	height int
}

func NewTexture(img Image) *Texture {
	id := gl.CreateTexture()

	gl.BindTexture(gl.TEXTURE_2D, id)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

	if img.Data() == nil {
		panic("Texture image data is nil.")
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.RGBA, gl.UNSIGNED_BYTE, img.Data())

	return &Texture{id, img.Width(), img.Height()}
}

// Width returns the width of the texture.
func (t *Texture) Width() float32 {
	return float32(t.width)
}

// Height returns the height of the texture.
func (t *Texture) Height() float32 {
	return float32(t.height)
}

func (t *Texture) Texture() *webgl.Texture {
	return t.id
}

func (r *Texture) View() (float32, float32, float32, float32) {
	return 0.0, 0.0, 1.0, 1.0
}

type Sound struct {
	*audio.Player
}

type Sprite struct {
	Position Point
	Scale    Point
	Anchor   Point
	Rotation float32
	Color    uint32
	Alpha    float32
	Region   *Region
}

func NewSprite(region *Region, x, y, scale float32) Sprite {
	return Sprite{
		Position: Point{x, y},
		Scale:    Point{scale, scale},
		Anchor:   Point{0, 0},
		Rotation: 0,
		Color:    0xffffff,
		Alpha:    1,
		Region:   region,
	}
}

func (s *Sprite) Render(batch *Batch, cam *Camera) {
	var position = cam.TranslateToScreen(s.Position)
	batch.Draw(s.Region, position.X, position.Y, s.Anchor.X, s.Anchor.Y, s.Scale.X*cam.Zoom, s.Scale.Y*cam.Zoom, s.Rotation, s.Color, s.Alpha)
}

func (s *Sprite) Width() float32 {
	return s.Region.width * s.Scale.X
}

func (s *Sprite) Height() float32 {
	return s.Region.height * s.Scale.Y
}

const batchVert = ` 
attribute vec2 in_Position;
attribute vec4 in_Color;
attribute vec2 in_TexCoords;

uniform vec2 uf_Projection;

varying vec4 var_Color;
varying vec2 var_TexCoords;

const vec2 center = vec2(-1.0, 1.0);

void main() {
  var_Color = in_Color;
  var_TexCoords = in_TexCoords;
	gl_Position = vec4(in_Position.x / uf_Projection.x + center.x,
										 in_Position.y / -uf_Projection.y + center.y,
										 0.0, 1.0);
}`

const batchFrag = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;
varying vec2 var_TexCoords;

uniform sampler2D uf_Texture;

void main (void) {
  gl_FragColor = var_Color * texture2D(uf_Texture, var_TexCoords);
}`

const batchLineVert = ` 
attribute vec2 in_Position;
attribute vec4 in_Color;

uniform vec2 uf_Projection;

varying vec4 var_Color;

const vec2 center = vec2(-1.0, 1.0);

void main() {
  var_Color = in_Color;
	gl_Position = vec4(in_Position.x / uf_Projection.x + center.x,
										 in_Position.y / -uf_Projection.y + center.y,
										 0.0, 1.0);
}`

const batchLineFrag = `
#ifdef GL_ES
#define LOWP lowp
precision mediump float;
#else
#define LOWP
#endif

varying vec4 var_Color;

void main (void) {
  gl_FragColor = var_Color;
}`
