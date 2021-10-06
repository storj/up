package main

import (
	sjr "github.com/elek/sjr/pkg"
	"log"
)

func main() {
	err := sjr.RootCmd.Execute()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}
