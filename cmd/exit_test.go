package cmd

import (
	"errors"
	"fmt"
	"testing"
)

func TestExitCodeFrom(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, 0},
		{"ErrNotFound", ErrNotFound, 2},
		{"ErrIOOrNetwork", ErrIOOrNetwork, 3},
		{"wrapped ErrNotFound", fmt.Errorf("license not found: MIT: %w", ErrNotFound), 2},
		{"wrapped ErrIOOrNetwork", fmt.Errorf("simulated: %w", ErrIOOrNetwork), 3},
		{"plain error", errors.New("usage: wrong args"), 1},
		{"other wrapped", fmt.Errorf("outer: %w", errors.New("inner")), 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := exitCodeFrom(tt.err)
			if got != tt.want {
				t.Errorf("exitCodeFrom() = %d, want %d", got, tt.want)
			}
		})
	}
}
