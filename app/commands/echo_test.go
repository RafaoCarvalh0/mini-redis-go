package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleEcho(t *testing.T) {
	type args struct {
		args   []string
		store  *map[string]Entry
		config server_config.ServerConfig
	}
	tests := []struct {
		name                string
		args                args
		want                string
		checkCharactesCount func(t *testing.T, args []string)
	}{
		{
			name: "returns provided ECHO message in RESP2 format",
			args: args{
				args:   []string{"ECHO", "foo"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "$3\r\nfoo\r\n",
		},
		{
			name: "returns provided ECHO message with the correct count of characters after $ sign",
			args: args{
				args:   []string{"ECHO", "foobar"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "$6\r\nfoobar\r\n",
			checkCharactesCount: func(t *testing.T, args []string) {
				got := handleEcho(args, &map[string]Entry{}, server_config.ServerConfig{})
				if len(got) < 2 || got[0] != '$' {
					t.Errorf("unexpected format: %q", got)
					return
				}

				endIdx := 1
				for endIdx < len(got) && got[endIdx] != '\r' {
					endIdx++
				}

				lengthStr := got[1:endIdx]
				if lengthStr != "6" {
					t.Errorf("expected length '6', got '%s'", lengthStr)
				}
			},
		},
		{
			name: "returns wrong number of arguments error when more than 2 arguments are provided",
			args: args{
				args:   []string{"ECHO", "foo", "bar"},
				store:  &map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'echo' command\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handleEcho(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handleEcho() = %v, want %v", got, tt.want)
			}

			if tt.checkCharactesCount != nil {
				tt.checkCharactesCount(t, tt.args.args)
			}
		})
	}
}
