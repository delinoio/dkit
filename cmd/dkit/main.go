package main

import (
	"os"

	"github.com/delinoio/dkit/internal/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
