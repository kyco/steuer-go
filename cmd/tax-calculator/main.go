package main

import (
	"fmt"
	"os"

	"tax-calculator/internal/ui"
)

func main() {
	if err := ui.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}