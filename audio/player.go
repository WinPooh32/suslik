package audio

import (
	"sync"

	"github.com/hajimehoshi/oto/v2"
)

type Player struct {
	once sync.Once

	player oto.Player
	sb     *soundbuf

	pause bool
	loop  bool
}

func (p *Player) Play(loop bool) {
	p.loop = loop
	p.sb.Loop(loop)

	if !p.player.IsPlaying() {
		if p.sb.finished {
			p.Seek(0)
		} else {
			p.player.Play()
		}
	} else {
		p.Seek(0)
	}
}

func (p *Player) Pause() {
	if !p.player.IsPlaying() {
		return
	}
	p.player.Pause()
}

func (p *Player) Seek(ms int) {
	p.player.Reset()
	p.sb.Seek(ms)
	p.player.Play()
}

func (p *Player) Playing() bool {
	return p.player.IsPlaying()
}

func (p *Player) Close() error {
	var err error

	p.once.Do(func() {
		err = p.player.Close()
	})

	return err
}

func NewPlayer(data []byte) *Player {
	sb := newSoundbuf(data)

	var p = &Player{
		sb:     sb,
		player: otoContext.NewPlayer(sb),

		once: sync.Once{},

		pause: false,
		loop:  false,
	}

	return p
}
