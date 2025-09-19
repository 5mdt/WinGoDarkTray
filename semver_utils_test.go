package main

import (
	"reflect"
	"testing"
)

func TestPadVersionParts(t *testing.T) {
	tests := []struct {
		name      string
		parts     []string
		targetLen int
		want      []string
	}{
		{
			name:      "pad from 2 to 3",
			parts:     []string{"1", "2"},
			targetLen: 3,
			want:      []string{"1", "2", "0"},
		},
		{
			name:      "pad from 1 to 3",
			parts:     []string{"1"},
			targetLen: 3,
			want:      []string{"1", "0", "0"},
		},
		{
			name:      "no padding needed",
			parts:     []string{"1", "2", "3"},
			targetLen: 3,
			want:      []string{"1", "2", "3"},
		},
		{
			name:      "already longer than target",
			parts:     []string{"1", "2", "3", "4"},
			targetLen: 3,
			want:      []string{"1", "2", "3", "4"},
		},
		{
			name:      "empty parts",
			parts:     []string{},
			targetLen: 2,
			want:      []string{"0", "0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := padVersionParts(tt.parts, tt.targetLen)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("padVersionParts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "a greater than b",
			a:    5,
			b:    3,
			want: 5,
		},
		{
			name: "b greater than a",
			a:    2,
			b:    7,
			want: 7,
		},
		{
			name: "a equals b",
			a:    4,
			b:    4,
			want: 4,
		},
		{
			name: "negative numbers",
			a:    -2,
			b:    -5,
			want: -2,
		},
		{
			name: "zero and positive",
			a:    0,
			b:    1,
			want: 1,
		},
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
