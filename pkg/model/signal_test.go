package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	tt := []struct {
		name          string
		input         []byte
		expectedType  rune
		expectedValue float64
	}{
		{
			name:          "EmptyInput",
			input:         []byte{},
			expectedType:  0,
			expectedValue: 0,
		},
		{
			name:          "A_4252",
			input:         []byte{65, 16, 156, 237, 13, 0, 0, 0},
			expectedType:  65,
			expectedValue: 4252,
		},
		{
			name:          "S_0_case1",
			input:         []byte{83, 0, 0, 83, 13, 0, 0, 0},
			expectedType:  83,
			expectedValue: 0,
		},
		{
			name:          "P_1373",
			input:         []byte{80, 5, 93, 178, 13, 0, 0, 0},
			expectedType:  80,
			expectedValue: 1373,
		},
		{
			name:          "B_4790_case1",
			input:         []byte{66, 18, 182, 10, 13, 0, 0, 0},
			expectedType:  66,
			expectedValue: 4790,
		},
		{
			name:          "A_4247",
			input:         []byte{65, 16, 151, 232, 13, 0, 0, 0},
			expectedType:  65,
			expectedValue: 4247,
		},
		{
			name:          "S_0_case2",
			input:         []byte{83, 0, 0, 83, 13, 0, 0, 0},
			expectedType:  83,
			expectedValue: 0,
		},
		{
			name:          "P_1380",
			input:         []byte{80, 5, 100, 185, 13, 0, 0, 0},
			expectedType:  80,
			expectedValue: 1380,
		},
		{
			name:          "B_4790_case2",
			input:         []byte{66, 18, 182, 10, 13, 0, 0, 0},
			expectedType:  66,
			expectedValue: 4790,
		},
	}

	t.Parallel()
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var actual = Signal(tc.input)
			assert.Equal(t, tc.expectedType, actual.Type())
			assert.Equal(t, tc.expectedValue, actual.Value())
		})
	}
}
