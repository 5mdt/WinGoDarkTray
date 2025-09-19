package main

import (
	"testing"
)

func TestNewApp(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "with version",
			version: "v1.2.3",
			want:    "v1.2.3",
		},
		{
			name:    "empty version defaults to v0.0.0",
			version: "",
			want:    "v0.0.0",
		},
		{
			name:    "with dev version",
			version: "v1.0.0-dev",
			want:    "v1.0.0-dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewApp(tt.version)
			if app == nil {
				t.Fatal("NewApp() returned nil")
			}
			if app.version != tt.want {
				t.Errorf("NewApp().version = %v, want %v", app.version, tt.want)
			}

			// Test that menu items are properly initialized as nil
			if app.toggleSystemItem != nil {
				t.Error("toggleSystemItem should be nil in new app")
			}
			if app.toggleAppItem != nil {
				t.Error("toggleAppItem should be nil in new app")
			}
			if app.toggleWindowsItem != nil {
				t.Error("toggleWindowsItem should be nil in new app")
			}
			if app.updateNowItem != nil {
				t.Error("updateNowItem should be nil in new app")
			}
		})
	}
}
