package main

import "testing"

func Test_lexer(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{name: "first test", want: "not implemented"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lexer(); got != tt.want {
				t.Errorf("lexer() = %v, want %v", got, tt.want)
			}
		})
	}
}
