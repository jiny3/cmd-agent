package tools

import (
	"reflect"
	"testing"
)

func Test_cmdExecutor(t *testing.T) {
	type args struct {
		args []any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "valid command with timeout",
			args: args{
				args: []any{"sleep 12", "3s"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "valid command with timeout",
			args: args{
				args: []any{"sleep 12", "20s"},
			},
			want:    []string{"", ""},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmdExecutor(tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("cmdExecutor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("cmdExecutor() = %v, want %v", got, tt.want)
			}
		})
	}
}
