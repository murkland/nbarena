package sound

import (
	"github.com/faiface/beep"
	"github.com/murkland/nbarena/bundle"
	"github.com/murkland/nbarena/state"
)

type playingSound struct {
	sound  *state.Sound
	stream beep.StreamSeeker
}

type scheduler struct {
	mixer            *beep.Mixer
	currentlyPlaying map[state.SoundID]playingSound
}

func (s *scheduler) Dispatch(b *bundle.Bundle, sounds map[state.SoundID]*state.Sound) {
	for id, cp := range s.currentlyPlaying {
		if cp.stream.Position() == cp.stream.Len() {
			delete(s.currentlyPlaying, id)
		}
	}

	for id, sound := range sounds {
		if cp, ok := s.currentlyPlaying[id]; ok {
			if cp.sound.Type == sound.Type && cp.sound.ElapsedTime == sound.ElapsedTime {
				continue
			}
			// Abort playback of the current sound.
			cp.stream.Seek(cp.stream.Len())
		}

		buf := b.Sounds[sound.Type]
		i := state.TicksToSampleOffset(buf.Format().SampleRate, sound.ElapsedTime)
		if i >= buf.Len() {
			continue
		}
		stream := buf.Streamer(i, buf.Len())
		s.mixer.Add(stream)
		s.currentlyPlaying[id] = playingSound{
			sound:  sound,
			stream: stream,
		}
	}
}

type Scheduler interface {
	Dispatch(b *bundle.Bundle, sounds map[state.SoundID]*state.Sound)
}

func NewScheduler(sr beep.SampleRate, mixer *beep.Mixer) Scheduler {
	return &scheduler{
		mixer:            mixer,
		currentlyPlaying: map[state.SoundID]playingSound{},
	}
}
