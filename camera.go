package suslik

type Camera struct {
	Position Point
	Viewport Point
	Zoom     float32
}

func (cam *Camera) TranslateToWorld(position Point) Point {
	position.X = position.X + cam.Position.X
	position.Y = position.Y + cam.Position.Y
	return position
}

func (cam *Camera) TranslateToScreen(position Point) Point {
	position.X = position.X - cam.Position.X
	position.Y = position.Y - cam.Position.Y
	return position
}

func (cam *Camera) MoveTo(position Point) {
	cam.Position = position
}

func NewCamera(position, viewport Point) *Camera {
	return &Camera{
		Position: Point{},
		Viewport: viewport,
		Zoom:     1.0,
	}
}
