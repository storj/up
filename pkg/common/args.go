// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package common

import (
	"fmt"
	"strings"
)

// ParseArgumentsWithSelector separate the selector and real arguments from os args.
func ParseArgumentsWithSelector(args []string, argNo int) ([]string, []string, error) {
	if argNo > len(args) {
		return nil, nil, fmt.Errorf("not enough arguments (required <selector> + %d)", argNo)
	}
	realArgs := args[len(args)-argNo:]

	selectors := make([]string, 0)
	for i := 0; i < len(args)-argNo; i++ {
		selectors = append(selectors, strings.Split(args[i], ",")...)
	}

	return selectors, realArgs, nil
}

// SplitArgsSelector1 splits selector from arguments leaving one extra separate argument.
// This is useful for splitting "[<selector>...] arg" or "<selector>... arg".
func SplitArgsSelector1(args []string) ([]string, string) {
	return args[:len(args)-1], args[len(args)-1]
}
