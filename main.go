package main

import (
	"darkroom/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err) // your colored message
		os.Exit(1)
	}
}
