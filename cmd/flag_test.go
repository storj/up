// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package cmd

import (
	"reflect"
	"testing"
)

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
			if got := removeFlag(tt.command, tt.param); !reflect.DeepEqual(got, tt.want) {
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
			if got := setFlag(tt.command, tt.param); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}
