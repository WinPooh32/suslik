package audio

import (
	"bytes"
	"io"
	"sync"

	"github.com/hajimehoshi/oto"
)

type Player struct {
	ch     chan int
	data   *bytes.Reader
	player *oto.Player

	once sync.Once
	done chan struct{}

	playing bool
	loop    bool

	start chan struct{}
	pause chan struct{}
}

func (p *Player) Play(loop bool) {
	if p.playing {
		return
	}
	p.playing = true
	p.loop = loop
	p.start <- struct{}{}
}

func (p *Player) Pause() {
	if !p.playing {
		return
	}
	p.playing = false
	p.pause <- struct{}{}
}

func (p *Player) Seek(ms int) {
	p.ch <- ms
}

func (p *Player) Playing() bool {
	return p.playing
}

func (p *Player) Close() error {
	var err error

	p.once.Do(func() {
		close(p.done)
		close(p.start)
		close(p.pause)
		err = p.player.Close()
	})

	return err
}

func NewPlayer(data []byte) *Player {

	var p = &Player{
		ch:     make(chan int, 1),
		data:   bytes.NewReader(data),
		player: otoContext.NewPlayer(),

		once: sync.Once{},
		done: make(chan struct{}),

		playing: false,
		loop:    false,

		start: make(chan struct{}, 1),
		pause: make(chan struct{}, 1),
	}

	// p.pause <- struct{}{}

	go func() {
		var playing bool

		for {

			select {
			case <-p.pause:
				// <-p.start
				playing = false
			case <-p.start:
				playing = true
			default:
			}

			select {
			case <-p.done:
				break
			case ms := <-p.ch:
				// TODO convert ms to buffer's offset.
				p.data.Seek(int64(ms), io.SeekStart)
			default:
			}

			var b [bufferSize]byte

			// Always write zeros, otherwise it locks on creating new player.
			// issue: https://github.com/hajimehoshi/oto/issues/117
			if playing {
				n, err := p.data.Read(b[:])
				if err == io.EOF && p.loop {
					p.data.Seek(0, io.SeekStart)
					p.data.Read(b[n:])
				}
			}

			p.player.Write(b[:])
		}
	}()

	return p
}
