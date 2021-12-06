package addon

import (
	"path"
	"sort"
	"strings"
	"time"

	"github.com/WinPooh32/suslik"
)

var DefaultFrameTime = 100 * time.Millisecond

type AnimationFrame struct {
	AnimationFrameSize `json:"frame"`
}

type AnimationFrameSize struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type ActInfo struct {
	FrameTime int `json:"frame_time"`
}

type AnimationInfo struct {
	Frames map[string]AnimationFrame `json:"frames"`
	Acts   map[string]ActInfo        `json:"acts"`
}

func (info *AnimationInfo) Regions(texture *suslik.Texture, prefix string) []*suslik.Region {
	var regions []*suslik.Region
	var list []string

	for name := range info.Frames {
		if strings.HasPrefix(name, prefix) {
			list = append(list, name)
		}
	}

	sort.Strings(list)

	for _, name := range list {
		frame := info.Frames[name]
		regions = append(regions, suslik.NewRegion(texture, frame.X, frame.Y, frame.W, frame.H))
	}

	return regions
}

type AnimationAct struct {
	control *suslik.Animation
	frames  []*suslik.Region
}

type Animator struct {
	acts    map[string]AnimationAct
	texture *suslik.Texture
}

func (animator *Animator) Draw(batch *suslik.Batch, name string, position, origin, scale suslik.Point, rot, alpha float32, color uint32) {
	if act, ok := animator.acts[name]; ok {
		batch.Draw(
			act.frames[act.control.NextFrame()],
			position.X, position.Y,
			origin.X, origin.Y,
			scale.X, scale.Y,
			rot,
			color,
			alpha,
		)
	}
}

func (animator *Animator) Playing(name string) bool {
	if _, ok := animator.acts[name]; ok {
		return animator.acts[name].control.Playing
	}
	return false
}

func (animator *Animator) Play(name string, frame int) {
	if _, ok := animator.acts[name]; ok {
		animator.acts[name].control.Play(frame)
	}
}

func (animator *Animator) Stop(name string) {
	if _, ok := animator.acts[name]; ok {
		animator.acts[name].control.Stop()
	}
}

func (animator *Animator) SetFrameTime(name string, dur time.Duration) {
	if _, ok := animator.acts[name]; ok {
		animator.acts[name].control.SetDelay(dur)
	}
}

func (animator *Animator) SetFrameTimeAll(dur time.Duration) {
	for name := range animator.acts {
		animator.acts[name].control.SetDelay(dur)
	}
}

func MakeAnimator(texture *suslik.Texture, info AnimationInfo) Animator {
	var acts []string

	for name := range info.Frames {
		var dir = path.Dir(name)
		if dir != "." {
			acts = append(acts, dir)
		}
	}

	var animator = Animator{
		acts:    make(map[string]AnimationAct, len(acts)),
		texture: texture,
	}

	if len(acts) == 0 {
		var name = ""
		animator.acts[name] = makeAnimationAct(name, texture, info)
	} else {
		for _, name := range acts {
			animator.acts[name] = makeAnimationAct(name, texture, info)
		}
	}

	return animator
}

func makeAnimationAct(name string, texture *suslik.Texture, info AnimationInfo) AnimationAct {
	var frames = info.Regions(texture, name)
	var actInfo, _ = info.Acts[name]
	var frameTime = time.Duration(actInfo.FrameTime) * time.Millisecond

	if frameTime <= 0 {
		frameTime = DefaultFrameTime
	}

	return AnimationAct{
		control: suslik.NewAnimation(len(frames), frameTime),
		frames:  frames,
	}
}
