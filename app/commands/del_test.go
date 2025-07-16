package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleDel(t *testing.T) {
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
			name: "delete one existing key",
			args: args{
				args: []string{"DEL", "key1"},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: ":1\r\n",
		},
		{
			name: "delete three existing keys",
			args: args{
				args: []string{"DEL", "key1", "key2", "key3"},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
					"key2": {Value: "value2", ExpiryTime: 0},
					"key3": {Value: "value3", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: ":3\r\n",
		},
		{
			name: "delete three keys where one doesn't exist",
			args: args{
				args: []string{"DEL", "key1", "key2", "nonexistent"},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
					"key2": {Value: "value2", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: ":2\r\n",
		},
		{
			name: "delete without arguments",
			args: args{
				args:   []string{"DEL"},
				store:  map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'del' command\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := handleDel(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handleDel() = %v, want %v", got, tt.want)
			}
		})
	}
}
