package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleEcho(t *testing.T) {
	type args struct {
		args   []string
		store  map[string]Entry
		config server_config.ServerConfig
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns provided ECHO message in RESP2 format",
			args: args{
				args:   []string{"ECHO", "foo"},
				store:  map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "$3\r\nfoo\r\n",
		},
		{
			name: "returns wrong number of arguments error when more than 2 arguments are provided",
			args: args{
				args:   []string{"ECHO", "foo", "bar"},
				store:  map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'echo' command\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := handleEcho(tt.args.args, tt.args.store, tt.args.config)
			if got != tt.want {
				t.Errorf("handleEcho() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleEcho_checkCharactesCount(t *testing.T) {
	t.Parallel()

	args := []string{"ECHO", "foobar"}

	echoResponse := handleEcho(args, map[string]Entry{}, server_config.ServerConfig{})

	if len(echoResponse) < 2 || echoResponse[0] != '$' {
		t.Errorf("unexpected format: %q", echoResponse)
		return
	}

	endIdx := 1
	for endIdx < len(echoResponse) && echoResponse[endIdx] != '\r' {
		endIdx++
	}

	lengthStr := echoResponse[1:endIdx]
	if lengthStr != "6" {
		t.Errorf("expected length '6', got '%s'", lengthStr)
	}
}
