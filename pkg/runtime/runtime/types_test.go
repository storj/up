// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package runtime

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServiceInstanceFromIndexedName(t *testing.T) {
	require.Equal(t, ServiceInstance{"asd", 1}, ServiceInstanceFromIndexedName("asd2"))
	require.Equal(t, ServiceInstance{"asd", 0}, ServiceInstanceFromIndexedName("asd1"))
	require.Equal(t, ServiceInstance{"satellite-api", 1}, ServiceInstanceFromIndexedName("satellite-api2"))
	require.Equal(t, ServiceInstance{"satellite-api", 0}, ServiceInstanceFromIndexedName("satellite-api1"))
	require.Equal(t, ServiceInstance{"satellite-api", 0}, ServiceInstanceFromIndexedName("satellite-api"))
}

func Test_removeFlag(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		param   string
		want    []string
	}{
		{
			name:    "Simple argument removal",
			command: []string{"storagenode", "--default=dev"},
			param:   "default",
			want:    []string{"storagenode"},
		},
		{
			name:    "Argument removal with key value",
			command: []string{"storagenode", "--default=dev"},
			param:   "default=dev",
			want:    []string{"storagenode"},
		},
		{
			name:    "Argument removal attempt, param doesn't exist",
			command: []string{"storagenode", "--default=dev"},
			param:   "foobar",
			want:    []string{"storagenode", "--default=dev"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveFlag(tt.command, tt.param); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setFlag(t *testing.T) {
	tests := []struct {
		name    string
		command []string
		param   string
		want    []string
	}{
		{
			name:    "Add new arguments",
			command: []string{"storagenode", "--default=dev"},
			param:   "foobar=something",
			want:    []string{"storagenode", "--default=dev", "--foobar=something"},
		},
		{
			name:    "Change existing argument",
			command: []string{"storagenode", "--default=dev"},
			param:   "default=prod",
			want:    []string{"storagenode", "--default=prod"},
		},
		{
			name:    "Change existing argument with one -",
			command: []string{"storagenode", "-default=dev"},
			param:   "default=prod",
			want:    []string{"storagenode", "-default=prod"},
		},
		{
			name:    "Add boolean flag",
			command: []string{"storagenode", "--default=dev"},
			param:   "foobar",
			want:    []string{"storagenode", "--default=dev", "--foobar"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetFlag(tt.command, tt.param); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
