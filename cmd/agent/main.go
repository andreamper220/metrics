package main

import (
	"github.com/andreamper220/metrics.git/internal/agent"
)

func main() {
	agent.ParseFlags()
	if err := agent.Run(); err != nil {
		panic(err)
	}
}
