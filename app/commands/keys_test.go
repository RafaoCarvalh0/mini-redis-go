package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleKeys(t *testing.T) {
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
			name: "returns error when not enough arguments are provided",
			args: args{
				args:   []string{"KEYS"},
				store:  map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'keys' command\r\n",
		},
		{
			name: "returns unknown command error when pattern is not *",
			args: args{
				args:   []string{"KEYS", "foo"},
				store:  map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}},
				config: server_config.ServerConfig{},
			},
			want: "-ERR unknown command foo\r\n",
		},
		{
			name: "returns empty array when store is empty and pattern is *",
			args: args{
				args:   []string{"KEYS", "*"},
				store:  map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "*0\r\n",
		},
		{
			name: "returns all keys in RESP2 array when pattern is *",
			args: args{
				args:   []string{"KEYS", "*"},
				store:  map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}, "baz": {Value: "qux", ExpiryTime: 0}},
				config: server_config.ServerConfig{},
			},
			want: "*2\r\n$3\r\nfoo\r\n$3\r\nbaz\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := handleKeys(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handleKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}
