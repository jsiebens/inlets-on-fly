package main

import (
	"github.com/jsiebens/inlets-on-fly/pkg/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
