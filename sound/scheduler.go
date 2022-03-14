package sound

import (
	"github.com/faiface/beep"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
)

type scheduler struct {
	mixer            *beep.Mixer
	currentlyPlaying map[state.SoundID]beep.StreamSeeker
}

func (s *scheduler) Dispatch(b *bundle.Bundle, sounds map[state.SoundID]*state.Sound) {
	for id, stream := range s.currentlyPlaying {
		if stream.Position() == stream.Len() {
			delete(s.currentlyPlaying, id)
		}
	}

	for id, sound := range sounds {
		if s.currentlyPlaying[id] != nil {
			continue
		}

		buf := b.Sounds[sound.Type]
		i := state.TicksToSampleOffset(buf.Format().SampleRate, sound.ElapsedTime)
		if i >= buf.Len() {
			continue
		}
		stream := buf.Streamer(i, buf.Len())
		s.mixer.Add(stream)
		s.currentlyPlaying[id] = stream
	}
}

type Scheduler interface {
	Dispatch(b *bundle.Bundle, sounds map[state.SoundID]*state.Sound)
}

func NewScheduler(sr beep.SampleRate, mixer *beep.Mixer) Scheduler {
	return &scheduler{
		mixer:            mixer,
		currentlyPlaying: map[state.SoundID]beep.StreamSeeker{},
	}
}
