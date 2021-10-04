package main

import sjr "github.com/elek/sjr/pkg"

func main() {
	err := sjr.RootCmd.Execute()
	if err != nil {
		panic(err)
	}
}
