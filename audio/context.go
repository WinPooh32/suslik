package audio

import (
	"log"

	"github.com/hajimehoshi/oto"
)

var otoContext *oto.Context

const (
	sampleRate = 44100
	channelNum = 2
	bitDepth   = 2
	bufferSize = 8 << 10
)

func init() {
	var err error

	otoContext, err = oto.NewContext(sampleRate, channelNum, bitDepth, bufferSize)
	if err != nil {
		log.Fatalf("failed to start Oto context: %s\n", err)
	}
}
