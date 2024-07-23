package main

import (
	"github.com/andreamper220/metrics.git/internal/server"
)

func main() {
	server.ParseFlags()
	if err := server.Run(false); err != nil {
		panic(err)
	}
}
