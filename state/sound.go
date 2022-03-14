package state

import (
	"time"

	"github.com/faiface/beep"
	"github.com/murkland/nbarena/bundle"
)

type SoundID uint64

type Sound struct {
	id          SoundID
	ElapsedTime Ticks

	Type bundle.SoundType
}

func (s *Sound) ID() SoundID {
	return s.id
}

func (s *Sound) Clone() *Sound {
	return &Sound{
		s.id,
		s.ElapsedTime,
		s.Type,
	}
}

func (s *Sound) Step() {
	s.ElapsedTime++
}

func TicksToSampleOffset(sr beep.SampleRate, t Ticks) int {
	return sr.N(time.Duration(t) * time.Second / 60)
}
