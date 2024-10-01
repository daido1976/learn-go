package main

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output map[string]string
	}{
		{
			name: "Basic test with valid input",
			input: `go.mod: text/plain
go.sum: text/plain
main: application/x-mach-binary
main.go: text/x-c
main.py: text/plain
output: inode/directory
space txt: text/plain`,
			output: map[string]string{
				"go.mod":    "text/plain",
				"go.sum":    "text/plain",
				"main":      "application/x-mach-binary",
				"main.go":   "text/x-c",
				"main.py":   "text/plain",
				"output":    "inode/directory",
				"space txt": "text/plain",
			},
		},
		{
			name:   "Empty input",
			input:  ``,
			output: map[string]string{},
		},
		{
			name:  "Filename with spaces",
			input: `file with spaces.txt: text/plain`,
			output: map[string]string{
				"file with spaces.txt": "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseMimeType(tt.input)
			if !reflect.DeepEqual(result, tt.output) {
				t.Errorf("parse() = %v, want %v", result, tt.output)
			}
		})
	}
}
