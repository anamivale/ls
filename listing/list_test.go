package listing

import (
	"io/fs"
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
		{
			name: "Not a directory",
			args: args{
				path: ".mal", flags: options.Flags{
					All:       false,
					Long:      false,
					Time:      false,
					Recursive: false,
					Reverse:   false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDirContent(tt.args.path, tt.args.flags)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDirContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
