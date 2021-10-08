package common

import (
	"github.com/stretchr/testify/require"
	"testing"
)

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

func Test_FilterPrefixAndGroup(t *testing.T) {
	base := &SimplifiedCompose{
		Services: map[string]*ServiceConfig{
			"storagenode1": &ServiceConfig{
				Image: "foobar",
			},
			"storagenode2": &ServiceConfig{
				Image: "foobar",
			},
			"storagenode3": &ServiceConfig{
				Image: "foobar",
			},
			"satellite-api": &ServiceConfig{},
		},
	}

	res := base.FilterPrefixAndGroup("satellite-api", map[string][]string{})
	require.Equal(t, 1, len(res))

	res = base.FilterPrefixAndGroup("storagenode", map[string][]string{})
	require.Equal(t, 3, len(res))

	res = base.FilterPrefixAndGroup("all", map[string][]string{
		"all": []string{"storagenode", "satellite-api"},
	})
	require.Equal(t, 4, len(res))
}
