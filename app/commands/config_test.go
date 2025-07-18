package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleConfig(t *testing.T) {
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
			name: "returns error for wrong number of arguments",
			args: args{
				args:   []string{"CONFIG", "GET"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'config' command\r\n",
		},
		{
			name: "returns error for unknown subcommand",
			args: args{
				args:   []string{"CONFIG", "GET", "INVALID"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR unknown command 'INVALID'\r\n",
		},
		{
			name: "returns DIR configuration in RESP2 format",
			args: args{
				args:   []string{"CONFIG", "GET", "DIR"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{Dir: "/tmp/redis"},
			},
			want: "*2\r\n$3\r\ndir\r\n$10\r\n/tmp/redis\r\n",
		},
		{
			name: "returns DBFILENAME configuration in RESP2 format",
			args: args{
				args:   []string{"CONFIG", "GET", "DBFILENAME"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{DBFileName: "dump.rdb"},
			},
			want: "*2\r\n$9\r\ndbfilename\r\n$8\r\ndump.rdb\r\n",
		},
		{
			name: "returns error for unsupported subcommand",
			args: args{
				args:   []string{"CONFIG", "GET", "UNSUPPORTED"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR unknown command 'UNSUPPORTED'\r\n",
		},
		{
			name: "returns DIR configuration when Dir is empty",
			args: args{
				args:   []string{"CONFIG", "GET", "DIR"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{Dir: ""},
			},
			want: "*2\r\n$3\r\ndir\r\n$0\r\n\r\n",
		},
		{
			name: "returns DBFILENAME configuration when Dbfilename is empty",
			args: args{
				args:   []string{"CONFIG", "GET", "DBFILENAME"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{DBFileName: ""},
			},
			want: "*2\r\n$9\r\ndbfilename\r\n$0\r\n\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := handleConfig(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handleConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
