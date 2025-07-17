package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func TestHandleCommand(t *testing.T) {
	type args struct {
		args   []string
		store  *map[string]Entry
		config server_config.ServerConfig
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "returns PONG for PING command",
			args: args{
				args:   []string{"PING"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want:  "+PONG\r\n",
			want1: true,
		},
		{
			name: "returns value for GET command",
			args: args{
				args:   []string{"GET", "foo"},
				store:  &map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}},
				config: server_config.ServerConfig{},
			},
			want:  "$3\r\nbar\r\n",
			want1: true,
		},
		{
			name: "returns OK for SET command",
			args: args{
				args:   []string{"SET", "foo", "bar"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want:  "+OK\r\n",
			want1: true,
		},
		{
			name: "returns error for unknown command",
			args: args{
				args:   []string{"FOOBAR"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want:  "-ERR unknown command 'FOOBAR'\r\n",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := HandleCommand(tt.args.args, tt.args.store, tt.args.config)
			if got != tt.want {
				t.Errorf("HandleCommand() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("HandleCommand() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
