package main

import (
	"fmt"
	"os"

	"github.com/mhdi/shamsi-calendar/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
