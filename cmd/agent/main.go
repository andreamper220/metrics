package main

import (
	"fmt"
	"github.com/andreamper220/metrics.git/internal/agent"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\r\nBuild date: %s\r\nBuild commit: %s\r\n", buildVersion, buildDate, buildCommit)
	agent.ParseFlags()
	if err := agent.Run(nil, nil); err != nil {
		panic(err)
	}
}
