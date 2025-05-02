package main

import (
	"fmt"
	"os"

	"tax-calculator/internal/main/views"
)

func main() {
	if err := views.Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}