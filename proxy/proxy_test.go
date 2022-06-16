package proxy

import (
	"testing"
)

func TestProxy(t *testing.T) {
	tests := []struct {
		name    string
		backend string
		want    bool
	}{
		{"No protocoll", "backend:7000", false},
		{"with protocoll", "http://backend:7000", true},
		{"wrong uri", ";backend:7000", false},
		{"wrong format", "backend7000", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if err := recover(); err != nil && tt.want {
					t.Errorf("Unexpected panic with %v", tt.backend)
				}
			}()
			Backend = tt.backend
			if got, err := NewReverseProxy(); got.proxy == nil || err !=nil {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
