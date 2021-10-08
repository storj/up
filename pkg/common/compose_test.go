package common

import "testing"

func Test_hasIndexedPrefix(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "notprefix",
			args: args{
				s:      "storagenode1",
				prefix: "storagenodex",
			},
			want: false,
		},
		{
			name: "prefixed",
			args: args{
				s:      "storagenode1",
				prefix: "storagenode",
			},
			want: true,
		},
		{
			name: "prefix-but-more-than-index",
			args: args{
				s:      "storagenodex1",
				prefix: "storagenode",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIndexedPrefix(tt.args.s, tt.args.prefix); got != tt.want {
				t.Errorf("hasIndexedPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
