package main

import (
	"fmt"
	"net/url"
	"testing"
)

func Test_extractIdFrom(t *testing.T) {
	type args struct {
		url *url.URL
	}
	id := 1
	url, _ := url.Parse(fmt.Sprintf("https://example.com/%d", id))
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{url: url},
			want: id,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractIdFrom(tt.args.url); got != tt.want {
				t.Errorf("extractIdFrom() = %v, want %v", got, tt.want)
			}
		})
	}
}
