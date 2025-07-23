package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleGet(t *testing.T) {
	type args struct {
		args   []string
		store  map[string]Entry
		config server_config.ServerConfig
	}
	type testCase struct {
		name string
		args args
		want string
	}

	tests := []testCase{
		{
			name: "returns value in RESP2 format when key exists and is not expired",
			args: args{
				args:   []string{"GET", "foo"},
				store:  map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}},
				config: server_config.ServerConfig{},
			},
			want: "$3\r\nbar\r\n",
		},
		{
			name: "returns error when not enough arguments are provided",
			args: args{
				args:   []string{"GET"},
				store:  map[string]Entry{},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'get' command\r\n",
		},
		{
			name: "returns $-1 when key does not exist",
			args: args{
				args:   []string{"GET", "missing"},
				store:  map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}},
				config: server_config.ServerConfig{},
			},
			want: "$-1\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := handleGet(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handleGet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleGet_checkExpiredKeyWasDeletedFromStore(t *testing.T) {
	t.Parallel()

	expiredKey := "foo"

	args := []string{"GET", expiredKey}

	store := map[string]Entry{"foo": {Value: "bar", ExpiryTime: 1}}

	handleGet(args, store, server_config.ServerConfig{})

	if _, keyFound := store[expiredKey]; keyFound {
		t.Errorf("handleGet() = not deleting expired keys")
	}
}
