package audio

type Scale []Interval

func (s Scale) Shift(shift int) Scale {
	return Scale(append(s[shift:], s[:shift]...))
}

// Quantize returns a slice of quantized normalized voltage values representing the scale.
func (s Scale) Quantize(root float32, input float32) float32 {
	var i int
	for {
		if root >= input {
			break
		}
		root += s[i%len(s)].ToVoltage()
	}
	return root
}

const (
	A0 float32 = 0.034375
)

const (
	Tone     Interval = 2
	SemiTone Interval = 1
)

type Interval float32

func (i Interval) ToVoltage() float32 {
	return float32(1 / MaxAbsVoltage / 12)
}

const (
	voltsPerSemitone = 1.0
)

var (
	ChromaticScale = Scale{
		SemiTone,
	}
	MajorScale = Scale{
		Tone,
		Tone,
		SemiTone,
		Tone,
		Tone,
		Tone,
		SemiTone,
	}
	IonianMode     = MajorScale
	DorianMode     = MajorScale.Shift(1)
	PhrygianMode   = MajorScale.Shift(2)
	LydianMode     = MajorScale.Shift(3)
	MixolydianMode = MajorScale.Shift(4)
	AeolianMode    = MajorScale.Shift(5)
	LocrianMode    = MajorScale.Shift(6)
	MinorScale     = Scale{
		Tone,
		SemiTone,
		Tone,
		Tone,
		SemiTone,
		Tone,
		Tone,
	}
	MinorPentatonicScale = Scale{
		Tone,
		SemiTone,
		Tone,
		Tone,
		SemiTone,
		Tone,
		Tone,
	}
)

type Quantizer struct {
	Input Processor
	Scale
	tonic float32
}

func NewQuantizer(scale Scale) *Quantizer {
	return &Quantizer{
		Scale: scale,
	}
}

func (q *Quantizer) Process(data []float32, channels int) {
	if q.Input != nil {
		q.Input.Process(data, channels)
	}
	var lastValue float32
	var lastQuantizedValue float32
	for i := 0; i < len(data); i += channels {
		for ch := 0; ch < channels; ch++ {
			if lastValue == data[ch+i] {
				data[ch+i] = lastQuantizedValue
			} else {
				lastQuantizedValue = q.Scale.Quantize(A0, data[ch+i])
				lastValue = data[ch+i]
				data[ch+i] = lastQuantizedValue
			}
		}
	}
}
