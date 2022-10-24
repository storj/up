// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package main

import (
	"storj.io/storj-up/cmd"
	_ "storj.io/storj-up/cmd/build"
	_ "storj.io/storj-up/cmd/container"
	_ "storj.io/storj-up/cmd/history"
	_ "storj.io/storj-up/cmd/modify"
)

func main() {
	cmd.Execute()
}
