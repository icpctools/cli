package commands

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKotlinBaseEntryPoint(t *testing.T) {
	testcases := []struct {
		input          string
		expectedOutput string
	}{
		{"main", "Main"},
		{"main append", "Main_append"},
		{"mainClass", "MainClass"},
		{"3main", "_3main"},
		{"#main#", "__main_"},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			output := kotlinBaseEntryPoint(tc.input)
			assert.EqualValues(t, tc.expectedOutput, output)
		})
	}
}
