package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handlePing(t *testing.T) {
	type args struct {
		args   []string
		store  *map[string]Entry
		config server_config.ServerConfig
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns PONG in RESP2 format when PING is provided",
			args: args{
				args:   []string{"PING"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "+PONG\r\n",
		},
		{
			name: "returns PONG in RESP2 format when PING and more arguments are provided",
			args: args{
				args:   []string{"PING", "PONG", "foo", "bar"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "+PONG\r\n",
		},
		{
			name: "returns an empty string when a different command is provided",
			args: args{
				args:   []string{"KEYS", "*"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "",
		},
		{
			name: "returns an empty string when no argument is provided",
			args: args{
				args:   []string{""},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := handlePing(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handlePing() = %v, want %v", got, tt.want)
			}
		})
	}
}
