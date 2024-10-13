package main

import (
	"fmt"
	"os"

	"github.com/inferablehq/sentinel/pkg/strategies"
)

func getConfig() strategies.Config {
	strategy := os.Getenv("STRATEGY")

	fmt.Println("Using strategy:", strategy)

	switch strategy {
	case "e2e":
		return strategies.E2EConfig()
	default:
		return strategies.MinimalConfig()
	}
}
