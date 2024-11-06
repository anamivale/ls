package listing

import (
	"io/fs"
	"reflect"
	"testing"

	"github.com/anamivale/ls/options"
)

func TestGetDirContent(t *testing.T) {
	type args struct {
		path  string
		flags options.Flags
	}
	tests := []struct {
		name    string
		args    args
		want    []fs.DirEntry
		wantErr bool
	}{
		{
			name: "Current dir",
			args: args{
				path: ".", flags: options.Flags{
					All:       false,
					Long:      false,
					Time:      false,
					Recursive: false,
					Reverse:   false,
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDirContent(tt.args.path, tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDirContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDirContent() = %v, want %v", got, tt.want)
			}
		})
	}
}
