// Copyright (C) 2022 Storj Labs, Inc.
// See LICENSE for copying information.

package standalone

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddConfig(t *testing.T) {
	s := service{
		render: func(s string) (string, error) {
			return s, nil
		},
		config: []string{
			"# address to listen on for debug endpoints",
			"# debug.addr: 127.0.0.1:0",
			"",
			"# expose control panel",
			"# debug.control: true",
			"",
			"# If set, a path to write a process trace SVG to",
			"# debug.trace-out: \"\"",
		},
	}
	require.NoError(t, s.AddConfig("STORJ_DEBUG_CONTROL", "false"))
	require.Equal(t, "debug.control: false", s.config[4])

	require.NoError(t, s.AddConfig("STORJ_DEBUG_TRACE_OUT", "xxx"))
	require.Equal(t, "debug.trace-out: \"xxx\"", s.config[7])

	require.NoError(t, s.AddConfig("new", "xxx"))
	require.Equal(t, "new: xxx", s.config[8])
}

func TestCamelToUpperCase(t *testing.T) {
	require.Equal(t, "STORJ_DEBUG_CONTROL", camelToUpperCase("debug.control"))
}
