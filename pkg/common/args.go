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
