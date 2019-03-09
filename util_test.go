package audio

import (
	"testing"
)

func TestNormalisedCVToFrequency(t *testing.T) {

	testCases := []struct {
		normalizedVoltage float32
		expectedFrequency float32
	}{
		{
			0.55,
			440,
		},
		{
			1,
			2093,
		},
		{
			-0.2,
			32.7,
		},
	}

	for _, tc := range testCases {
		if freq := NormalisedCVToFrequency(tc.normalizedVoltage); freq != tc.expectedFrequency {
			t.Errorf("Test failed, expected %v, got %v", tc.expectedFrequency, freq)
		}
	}

}
