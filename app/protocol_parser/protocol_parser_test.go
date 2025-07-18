package protocol_parser

import (
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestGetRESP2ArgsFromConn(t *testing.T) {
	type args struct {
		conn net.Conn
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr error
	}{
		{
			name: "returns a received RESP2 command parsed into a slice",
			args: args{
				conn: buildConnWithRESP2command("*2\r\n$4\r\nKEYS\r\n$1\r\n*\r\n"),
			},
			want:    []string{"KEYS", "*"},
			wantErr: nil,
		},
		{
			name: "returns a error for invalid RESP2 command",
			args: args{
				conn: buildConnWithRESP2command("2\r\n$4\r\nKEYS\r\n$1\r\n*\r\n"),
			},
			want:    []string{},
			wantErr: fmt.Errorf("invalid resp2 command"),
		},
		{
			name: "returns error for incorrect length on RESP2 command",
			args: args{
				conn: buildConnWithRESP2command("*2\r\n$4\r\nKEYS\r\n$1\r\n*\r\n$1\r\nfoo\r\n"),
			},
			want:    []string{},
			wantErr: fmt.Errorf("invalid resp2 command"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := GetRESP2ArgsFromConn(tt.args.conn)
			if tt.wantErr != nil {
				if err == nil {
					t.Errorf("GetRESP2ArgsFromConn() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.wantErr.Error() {
					t.Errorf("GetRESP2ArgsFromConn() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRESP2ArgsFromConn() = %v, want %v", got, tt.want)
			}
		})
	}
}

func buildConnWithRESP2command(resp2Command string) net.Conn {
	pr, pw := net.Pipe()
	go func() {
		pw.Write([]byte(resp2Command))
		pw.Close()
	}()
	return pr
}
