package commands

import (
	"bytes"
	"mini-redis-go/app/server_config"
	"os"
	"testing"
)

func Test_handleSave(t *testing.T) {
	type args struct {
		args   []string
		store  map[string]Entry
		config server_config.ServerConfig
	}
	type testCase struct {
		name      string
		args      args
		want      string
		checkFile func(t *testing.T, config server_config.ServerConfig)
	}

	tests := []testCase{
		{
			name: "returns error when directory or dbfilename is missing",
			args: args{
				args:   []string{"SAVE"},
				store:  map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}},
				config: server_config.ServerConfig{Dir: "", DBFileName: ""},
			},
			want: "-ERR no directory or dbfilename provided\r\n",
		},
		{
			name: "returns OK and saved file contains header and version",
			args: args{
				args:   []string{"SAVE"},
				store:  map[string]Entry{"foo": {Value: "bar", ExpiryTime: 0}},
				config: server_config.ServerConfig{},
			},
			want: "+OK\r\n",
			checkFile: func(t *testing.T, config server_config.ServerConfig) {
				filePath := config.Dir + "/dump.rdb"
				data, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("expected file to be created, got error: %v", err)
				}
				if len(data) == 0 {
					t.Errorf("expected file to have content, but it is empty")
				}
				if !bytes.HasPrefix(data, []byte("REDIS0001")) {
					t.Errorf("expected file to start with REDIS0001, got: %v", data[:8])
				}
				if !bytes.Contains(data, []byte("foo")) {
					t.Errorf("expected file to contain key 'foo'")
				}
				if !bytes.Contains(data, []byte("bar")) {
					t.Errorf("expected file to contain value 'bar'")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.checkFile != nil {
				tt.args.config.Dir = t.TempDir()
				tt.args.config.DBFileName = "dump.rdb"
			}
			got := handleSave(tt.args.args, tt.args.store, tt.args.config)
			if got != tt.want {
				t.Errorf("handleSave() = %v, want %v", got, tt.want)
			}
			if tt.checkFile != nil {
				tt.checkFile(t, tt.args.config)
			}
		})
	}
}
