package commands

import "testing"

func TestRedisError(t *testing.T) {
	type args struct {
		errMsg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns error with message",
			args: args{errMsg: "some error"},
			want: "-ERR some error\r\n",
		},
		{
			name: "returns unknown error with empty message",
			args: args{errMsg: ""},
			want: "-ERR unknown error\r\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RedisError(tt.args.errMsg); got != tt.want {
				t.Errorf("RedisError() = %v, want %v", got, tt.want)
			}
		})
	}
}
