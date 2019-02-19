package audio

import (
	"fmt"

	"github.com/rakyll/portmidi"
)

const (
	MidiNoteOn = 144
	Control    = 176
)

type MidiInput struct {
	*portmidi.Stream
	Gate    *Constant
	CV      *Constant
	control map[int64]*Constant
}

func NewMidiInput(stream *portmidi.Stream) *MidiInput {
	m := &MidiInput{
		Stream:  stream,
		Gate:    NewConstant(0).SetGlideMs(1),
		CV:      NewConstant(0).SetGlideMs(0),
		control: make(map[int64]*Constant),
	}
	go m.Listen()
	return m
}

func (m *MidiInput) Control(i int64) *Constant {
	m.control[i] = NewConstant(0).SetGlideMs(10)
	return m.control[i]
}

func (m *MidiInput) Listen() {
	var lastNote int64
	for msg := range m.Stream.Listen() {
		switch msg.Status {
		case MidiNoteOn:
			if lastNote != msg.Data1 && msg.Data2 == 0 {
				continue
			}
			if msg.Data2 != 0 {
				m.Gate.SetOffset(float32(msg.Data2) / 127)
			} else {
				m.Gate.SetOffset(0)
			}
			m.CV.SetOffset(MidiToNormalizedCV(msg.Data1))
			lastNote = msg.Data1
		case Control:
			fmt.Println(msg.Data1)
			if c, ok := m.control[msg.Data1]; ok {
				c.SetOffset((MidiToNormalizedCV(msg.Data2)))
			}
		}
	}
}
