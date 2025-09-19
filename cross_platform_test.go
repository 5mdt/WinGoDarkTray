// cross_platform_test.go

//go:build !windows
// +build !windows

package main

import (
	"testing"
)

// Cross-platform tests that don't depend on Windows APIs

func TestIsVersionNewerCrossPlatform(t *testing.T) {
	tests := []struct {
		latest  string
		current string
		want    bool
	}{
		{"0.9.9", "0.9.8", true},
		{"1.10.0", "1.2.3", true},
		{"1.2.3", "1.10.0", false},
		{"1.2.3", "1.2.3", false},
		{"1.2.3", "2.0.0", false},
		{"1.2.3", "invalid", false},
		{"1.2.3", "v1.2.3", false},
		{"1.2.4", "1.2.3", true},
		{"1.2", "1.2.0", false},
		{"1.3.0", "1.2.9", true},
		{"2.0.0", "1.9.9", true},
		{"invalid", "1.2.3", false},
		{"v1.10.0", "v1.2.3", true},
		{"v1.2.3", "1.2.3", false},
		{"v1.2.3", "v1.10.0", false},
		{"v1.2.4", "v1.2.3", true},
		{"v2.0.0", "v1.9.9", true},
	}

	for _, tt := range tests {
		t.Run(tt.latest+" vs "+tt.current, func(t *testing.T) {
			got := isVersionNewer(tt.latest, tt.current)
			if got != tt.want {
				t.Errorf("isVersionNewer(%q, %q) = %v; want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}

func TestMaxCrossPlatform(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{"a greater than b", 5, 3, 5},
		{"b greater than a", 2, 7, 7},
		{"a equals b", 4, 4, 4},
		{"negative numbers", -2, -5, -2},
		{"zero and positive", 0, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := max(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("max(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
