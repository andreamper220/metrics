package main

import (
	"github.com/andreamper220/metrics.git/internal/server/application"
)

func main() {
	application.ParseFlags()
	if err := application.Run(false); err != nil {
		panic(err)
	}
}
