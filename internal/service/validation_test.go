package service

import (
	"strings"
	"testing"
)

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{
			name:    "valid title",
			title:   "Learn Go",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
		},
		{
			name:    "spaces only",
			title:   "     ",
			wantErr: true,
		},
		{
			name:    "100 characters",
			title:   strings.Repeat("a", 100),
			wantErr: false,
		},
		{
			name:    "101 characters",
			title:   strings.Repeat("a", 101),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateTitle(tt.title)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateTitle(%q): got error %v, wantErr %v",
					tt.title,
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{
			name:        "valid description",
			description: "Study Backend Golang",
			wantErr:     false,
		},
		{
			name:        "empty description",
			description: "",
			wantErr:     false,
		},
		{
			name:        "500 characters",
			description: strings.Repeat("a", 500),
			wantErr:     false,
		},
		{
			name:        "501 characters",
			description: strings.Repeat("a", 501),
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDescription(tt.description)

			if (err != nil) != tt.wantErr {
				t.Errorf(
					"validateDescription(%q): got error %v, wantErr %v",
					tt.description,
					err,
					tt.wantErr,
				)
			}
		})
	}

}
