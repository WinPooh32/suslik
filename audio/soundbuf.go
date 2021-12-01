package audio

import (
	"fmt"
	"io"
)

type soundbuf struct {
	buf      []byte
	ptr      int
	loop     bool
	finished bool
}

func newSoundbuf(data []byte) *soundbuf {
	return &soundbuf{
		buf: data,
		ptr: 0,
	}
}

func (rb *soundbuf) Done() bool {
	return rb.finished
}

func (rb *soundbuf) Loop(enabled bool) {
	rb.loop = enabled
}

func (rb *soundbuf) Seek(pos int) int {
	if pos < len(rb.buf) {
		rb.ptr = pos
	} else {
		rb.ptr = len(rb.buf) - 1
	}
	rb.finished = false
	return rb.ptr
}

func (rb *soundbuf) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, fmt.Errorf("empty p")
	}

	ptr := rb.ptr

	if rb.loop {
		beg := ptr
		end := (ptr + len(p)) % len(rb.buf)

		if end > beg {
			copy(p, rb.buf[beg:end])
		} else {
			part := rb.buf[beg:]
			copy(p, part)

			if len(p) > len(part) {
				copy(p[len(part):], rb.buf[:end])
			}
		}

		rb.ptr = end

		return len(p), nil
	} else {
		n = copy(p, rb.buf[ptr:])

		if n == 0 {
			rb.finished = true
			return 0, io.EOF
		}

		rb.ptr = ptr + n

		return n, nil
	}
}
