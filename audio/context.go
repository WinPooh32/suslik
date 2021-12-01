package audio

import (
	"log"

	"github.com/hajimehoshi/oto/v2"
)

var otoContext *oto.Context

const (
	sampleRate = 44100
	channelNum = 2
	bitDepth   = 2
)

func init() {
	var err error
	var ready chan struct{}

	otoContext, ready, err = oto.NewContext(sampleRate, channelNum, bitDepth)
	if err != nil {
		log.Fatalf("failed to start Oto context: %s\n", err)
	}

	<-ready
}
