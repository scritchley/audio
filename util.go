package audio

import "math"

const (
	BaseFrequency float64 = 440
	BaseVoltage   float64 = 2.75
	MaxAbsVoltage float64 = 5
)

func NormalisedCVToFrequency(value float32) float32 {
	return float32(BaseFrequency / math.Pow(2, BaseVoltage) * math.Pow(2, float64(value)*MaxAbsVoltage))
}

func MidiToNormalizedCV(note int64) float32 {
	return (float32(note) / 60) - 1
}
