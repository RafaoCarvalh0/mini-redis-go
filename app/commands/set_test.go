package commands

import (
	"mini-redis-go/app/server_config"
	"testing"
)

func Test_handleSet(t *testing.T) {
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
			name: "stores a new key value pair",
			args: args{
				args: []string{"SET", "foo", "bar"},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: "+OK\r\n",
		},
		{
			name: "returns wrong number of arguments error when more than 3 arguments are provided",
			args: args{
				args: []string{"SET", "foo", "bar", "baz"},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'set' command\r\n",
		},
		{
			name: "returns wrong number of arguments error when less than 3 arguments are provided",
			args: args{
				args: []string{"SET", "foo"},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: "-ERR wrong number of arguments for 'set' command\r\n",
		},
		{
			name: "returns and empty string when an empty slice is provided",
			args: args{
				args: []string{""},
				store: map[string]Entry{
					"key1": {Value: "value1", ExpiryTime: 0},
				},
				config: server_config.ServerConfig{},
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := handleSet(tt.args.args, tt.args.store, tt.args.config); got != tt.want {
				t.Errorf("handleSet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleSet_StoreUpdate(t *testing.T) {
	store := map[string]Entry{
		"key1": {Value: "value1", ExpiryTime: 0},
	}
	config := server_config.ServerConfig{}
	args := []string{"SET", "foo", "bar"}

	handleSet(args, store, config)

	entry, ok := store["foo"]
	if !ok {
		t.Errorf("store should contain key 'foo'")
	} else if entry.Value != "bar" {
		t.Errorf("store[\"foo\"] = %v, want %v", entry.Value, "bar")
	}
}
