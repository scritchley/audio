package audio

type Delay struct {
	sampleRate   int
	delaySamples int
	buffer       []float32
}

func NewDelay() *Delay {
	return &Delay{
		sampleRate: DefaultSampleRate,
	}
}

func (d *Delay) SetDelay(delayMs float32) {
	d.delaySamples = int((delayMs * float32(d.sampleRate) / 1000))
}

func (d *Delay) Process(data []float32, channels int) {
	if len(d.buffer) < (d.delaySamples * channels) {
		d.buffer = append(d.buffer, make([]float32, (d.delaySamples*channels)-len(d.buffer))...)
	}
	for i := 0; i < len(data); i += channels {
		for ch := 0; ch < channels; ch++ {
			// data[ch+i] += d.buffer[(ch+i-1)%len(d.buffer)]
			d.buffer[(ch+i)%(d.delaySamples*channels)] += data[ch+i]
			data[ch+i] = d.buffer[(ch+i)%(d.delaySamples*channels)]
		}
	}
}
