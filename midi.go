package audio

import (
	"github.com/rakyll/portmidi"
)

const (
	MidiNoteOn  = 0x90
	MidiNoteOff = 0x80
)

type MidiInput struct {
	*portmidi.Stream
	Gate *Constant
	CV   *Constant
}

func NewMidiInput(stream *portmidi.Stream) *MidiInput {
	return &MidiInput{
		Stream: stream,
		Gate:   NewConstant(0),
		CV:     NewConstant(0),
	}
}

func (m *MidiInput) Listen() {
	for msg := range m.Stream.Listen() {
		switch msg.Data1 {
		case MidiNoteOn:
			m.Gate.SetOffset(1)
		case MidiNoteOff:
			m.Gate.SetOffset(0)
		}
	}
}
