// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseArgumentsWithSelector(t *testing.T) {
	type args struct {
		args  []string
		argNo int
	}
	tests := []struct {
		name         string
		args         args
		wantSelector []string
		wantArgs     []string
		wantErr      bool
	}{
		{
			name: "one selector and one arg",
			args: args{
				args:  []string{"satellite-api", "10"},
				argNo: 1,
			},
			wantSelector: []string{"satellite-api"},
			wantArgs:     []string{"10"},
		},
		{
			name: "one selector and two args",
			args: args{
				args:  []string{"satellite-api", "10", "20"},
				argNo: 2,
			},
			wantSelector: []string{"satellite-api"},
			wantArgs:     []string{"10", "20"},
		},
		{
			name: "comma separated selector",
			args: args{
				args:  []string{"satellite-api,storagenode", "10", "20"},
				argNo: 2,
			},
			wantSelector: []string{"satellite-api", "storagenode"},
			wantArgs:     []string{"10", "20"},
		},
		{
			name: "not enough arguments",
			args: args{
				args:  []string{"10"},
				argNo: 2,
			},
			wantErr: true,
		},
		{
			name: "zero arguments",
			args: args{
				args:  []string{"satellite-api"},
				argNo: 0,
			},
			wantSelector: []string{"satellite-api"},
			wantArgs:     []string{},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector, args, err := ParseArgumentsWithSelector(tt.args.args, tt.args.argNo)
			if tt.wantErr {
				require.NotNil(t, err)
				return
			}
			require.Equal(t, tt.wantSelector, selector)
			require.Equal(t, tt.wantArgs, args)
		})
	}
}
