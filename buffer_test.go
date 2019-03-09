package audio

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/gordonklaus/portaudio"
	"github.com/rakyll/portmidi"
)

func TestSelf(t *testing.T) {

	portaudio.Initialize()
	defer portaudio.Terminate()

	o := NewOscillator(Sawtooth)
	o.Frequency = NewConstant(0.55)

	stream, err := portaudio.OpenDefaultStream(0, 1, DefaultSampleRate, 0, func(out [][]float32) {
		o.Process(out[0], 1)
	})
	if err != nil {
		t.Fatal(err)
	}

	stream.Start()

	time.Sleep(time.Second * 10)

	stream.Stop()

}

func TestBuffer(t *testing.T) {

	fmt.Println(AeolianMode)

	portmidi.Initialize()
	defer portmidi.Terminate()

	in, err := portmidi.NewInputStream(portmidi.DefaultInputDeviceID(), 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	m := NewMidiInput(in)

	portaudio.Initialize()
	defer portaudio.Terminate()

	modA := NewOscillator(Sine)
	modA.Frequency = Add(m.CV, m.Control(84))
	modAAmp := NewAmplifier()
	modAAmp.Gain = m.Control(5)
	modAChain := NewProcessorChain(
		modA,
		modAAmp,
	)

	mod := NewOscillator(Sine)
	mod.Frequency = modAChain
	modAmp := NewAmplifier()
	modAmp.Gain = m.Control(71)
	modChain := NewProcessorChain(
		mod,
		modAmp,
	)

	o := NewOscillator(Sine)
	o.Frequency = m.CV
	o.Phase = modChain

	a := NewAmplifier()
	a.Gain = m.Gate

	f := NewFilter()
	f.CutOff = m.Control(72)
	f.Resonance = m.Control(73)

	// c := NewConstant(0.01)
	// c.SetGlideMs(5000)
	// q := NewQuantizer(MinorScale)
	// q.Input = c
	// o := NewOscillator(Sawtooth)
	// o.Frequency = q

	// cutOff := NewConstant(0.5)
	// cutOff.SetGlideMs(1000)
	// resonance := NewConstant(1)

	// f := NewFilter()
	// f.CutOff = cutOff
	// f.Resonance = resonance

	// d := NewDelay()
	// d.SetDelay(1000)

	stream, err := portaudio.OpenDefaultStream(0, 1, DefaultSampleRate, 0, func(out [][]float32) {
		o.Process(out[0], 1)
		f.Process(out[0], 1)
		a.Process(out[0], 1)
		// d.Process(out[0], 1)
	})
	if err != nil {
		t.Fatal(err)
	}

	stream.Start()

	// cutOff.SetOffset(0.0)

	// time.Sleep(time.Second)
	// c.SetOffset(0.1)

	// time.Sleep(5 * time.Second)

	// cutOff.SetOffset(1)

	// o.SetWaveFunc(Sawtooth)

	// time.Sleep(time.Second)
	// c.SetOffset(0.15)

	// time.Sleep(5 * time.Second)

	time.Sleep(time.Hour)

	stream.Stop()

}

func TestBasic(t *testing.T) {

	portaudio.Initialize()
	defer portaudio.Terminate()

	portmidi.Initialize()
	defer portmidi.Terminate()

	var deviceID, deviceOutID portmidi.DeviceID
	for i := 0; i < portmidi.CountDevices(); i++ {

		info := portmidi.Info(portmidi.DeviceID(i))
		if strings.Contains(info.Name, "Launchpad Mini") && info.IsInputAvailable {
			deviceID = portmidi.DeviceID(i)
		}
		if strings.Contains(info.Name, "Launchpad Mini") && info.IsOutputAvailable {
			deviceOutID = portmidi.DeviceID(i)
		}
	}

	out, err := portmidi.NewOutputStream(deviceOutID, 1024, 100)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			out.WriteShort(144, rand.Int63n(127), rand.Int63n(127))
			time.Sleep(16 * time.Millisecond)
		}
	}()

	in, err := portmidi.NewInputStream(deviceID, 1024)
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	m := NewMidiInput(in)

	mod := NewOscillator(Sine)
	mod.Frequency = m.CV

	o := NewOscillator(Sine)
	o.Frequency = m.CV
	o.Phase = mod

	a := NewAmplifier()
	a.Gain = m.Gate

	ch := NewProcessorChain(
		o,
		a,
	)

	stream, err := portaudio.OpenDefaultStream(0, 1, DefaultSampleRate, 0, func(out [][]float32) {
		ch.Process(out[0], 1)
	})
	if err != nil {
		t.Fatal(err)
	}

	stream.Start()

	time.Sleep(300 * time.Second)

	stream.Stop()

}
