package main

import (
	"fmt"
	"testing"
)

func Test_sendMsgText(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				text: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sendMsgText(tt.args.text)
		})
	}
}

func Test_a(t *testing.T) {
	a := fmt.Sprintf(""+
		"[ETH collision find address]\n"+
		"Total:%d\n"+
		"Speed:%d\n"+
		"Addrs: %d\n",
		18925693, 9932857, 0)
	print(a)
}
