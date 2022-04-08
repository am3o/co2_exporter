package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSignal(t *testing.T) {
	tt := []struct {
		input         []byte
		expectedType  rune
		expectedValue float64
	}{
		{
			input:         []byte{65, 16, 156, 237, 13, 0, 0, 0},
			expectedType:  65,
			expectedValue: 4252,
		},
		{
			input:         []byte{83, 0, 0, 83, 13, 0, 0, 0},
			expectedType:  83,
			expectedValue: 0,
		},
		{
			input:         []byte{80, 5, 93, 178, 13, 0, 0, 0},
			expectedType:  80,
			expectedValue: 1373,
		},
		{
			input:         []byte{66, 18, 182, 10, 13, 0, 0, 0},
			expectedType:  66,
			expectedValue: 4790,
		},
		{
			input:         []byte{65, 16, 151, 232, 13, 0, 0, 0},
			expectedType:  65,
			expectedValue: 4247,
		},
		{
			input:         []byte{83, 0, 0, 83, 13, 0, 0, 0},
			expectedType:  83,
			expectedValue: 0,
		},
		{
			input:         []byte{80, 5, 100, 185, 13, 0, 0, 0},
			expectedType:  80,
			expectedValue: 1380,
		},
		{
			input:         []byte{66, 18, 182, 10, 13, 0, 0, 0},
			expectedType:  66,
			expectedValue: 4790,
		},
	}

	t.Parallel()
	for _, tc := range tt {
		t.Run("", func(t *testing.T) {
			var actual = Signal(tc.input)
			assert.Equal(t, tc.expectedType, actual.Type())
			assert.Equal(t, tc.expectedValue, actual.Value())
		})
	}
}
