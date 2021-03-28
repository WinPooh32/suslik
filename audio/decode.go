package audio

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"github.com/WinPooh32/suslik/audio/decode/mp3"
	"github.com/WinPooh32/suslik/audio/decode/vorbis"
	"github.com/WinPooh32/suslik/audio/decode/wav"
)

var errNoDec = errors.New("decoder not found")

type closerSeekerReader interface {
	io.Reader
	io.Seeker
	io.Closer
}

// Decode decodes raw data according to file extension.
// Supported decoders: wav(.wav), vorbis(.ogg), mp3(.mp3).
func Decode(data []byte, ext string) ([]byte, error) {
	stream, err := decode(data, ext)
	if err != nil {
		return nil, err
	}
	decoded, err := ioutil.ReadAll(stream)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func decode(data []byte, ext string) (closerSeekerReader, error) {
	r := newReadSeekCloser(data)
	switch ext {
	case ".wav":
		return wav.Decode(r, sampleRate)
	case ".mp3":
		return mp3.Decode(r, sampleRate)
	case ".ogg":
		return vorbis.Decode(r, sampleRate)
	default:
		return nil, errNoDec
	}
}

type nopCloserSeekerReader struct {
	*bytes.Reader
}

func (r nopCloserSeekerReader) Close() error { return nil }

func newReadSeekCloser(data []byte) *nopCloserSeekerReader {
	return &nopCloserSeekerReader{bytes.NewReader(data)}
}
