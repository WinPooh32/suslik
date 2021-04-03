package suslik

import (
	"log"

	"github.com/WinPooh32/suslik/math"

	"github.com/ajhager/webgl"
)

type BatchLine struct {
	drawing      bool
	vertices     []float32
	vertexVBO    *webgl.Buffer
	index        int
	indices      []uint16
	indexVBO     *webgl.Buffer
	shader       *webgl.Program
	inPosition   int
	inColor      int
	ufProjection *webgl.UniformLocation
	projX        float32
	projY        float32
}

func NewBatchLine(width, height float32) *BatchLine {
	batch := new(BatchLine)

	batch.shader = LoadShader(batchLineVert, batchLineFrag)

	batch.inPosition = gl.GetAttribLocation(batch.shader, "in_Position")
	batch.inColor = gl.GetAttribLocation(batch.shader, "in_Color")
	batch.ufProjection = gl.GetUniformLocation(batch.shader, "uf_Projection")

	batch.vertices = make([]float32, 3*size)
	batch.indices = make([]uint16, 2*size)

	for i, j := 0, 0; i < size*2; i, j = i+2, j+2 {
		batch.indices[i+0] = uint16(j + 0)
		batch.indices[i+1] = uint16(j + 1)
	}

	batch.indexVBO = gl.CreateBuffer()
	batch.vertexVBO = gl.CreateBuffer()

	gl.EnableVertexAttribArray(batch.inPosition)
	gl.EnableVertexAttribArray(batch.inColor)

	batch.projX = width / 2
	batch.projY = height / 2

	return batch
}

func (b *BatchLine) Begin() {
	if b.drawing {
		log.Fatal("Batch.End() must be called first")
	}
	b.drawing = true

	gl.UseProgram(b.shader)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.indexVBO)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, b.indices, gl.STATIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, b.vertices, gl.DYNAMIC_DRAW)

	gl.VertexAttribPointer(b.inPosition, 2, gl.FLOAT, false, 12, 0)
	gl.VertexAttribPointer(b.inColor, 4, gl.UNSIGNED_BYTE, true, 12, 8)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
}

func (b *BatchLine) End() {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}
	if b.index > 0 {
		b.flush()
	}
	b.drawing = false
}

func (b *BatchLine) flush() {
	gl.Uniform2f(b.ufProjection, b.projX, b.projY)

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, b.vertices)
	gl.DrawElements(gl.LINES, 2*b.index, gl.UNSIGNED_SHORT, 0)

	b.index = 0
}

func (b *BatchLine) SetProjection(width, height float32) {
	b.projX = width / 2
	b.projY = height / 2
}

func (b *BatchLine) Draw(l Line, color uint32, transparency float32) {
	if !b.drawing {
		log.Fatal("Batch.Begin() must be called first")
	}

	red := (color >> 16) & 0xFF
	green := ((color >> 8) & 0xFF) << 8
	blue := (color & 0xFF) << 16
	alpha := uint32(transparency*255.0) << 24
	tint := math.Float32frombits((alpha | blue | green | red) & 0xfeffffff)

	idx := b.index * 6

	var s []float32
	if len(b.vertices) >= idx+6 {
		s = b.vertices[idx : idx+6]
	}

	if len(s) >= 6 {
		s[0] = l.P1.X
		s[1] = l.P1.Y
		s[2] = tint

		s[3] = l.P2.X
		s[4] = l.P2.Y
		s[5] = tint
	}

	b.index += 1

	if b.index >= size {
		b.flush()
	}
}
