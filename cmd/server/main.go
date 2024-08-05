package main

import (
	"fmt"
	"github.com/andreamper220/metrics.git/internal/server/application"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\r\nBuild date: %s\r\nBuild commit: %s\r\n", buildVersion, buildDate, buildCommit)
	application.ParseFlags()
	if err := application.Run(false); err != nil {
		panic(err)
	}
}
