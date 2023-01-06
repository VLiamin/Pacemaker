package main

import (
	"os"

	"yadro.com/m.konovalov/pcmkr/cmd/agents/ocf"
)

func main() { os.Exit(ocf.Run(Agent{}, os.Args[1:], ocf.OSEnvironment{}).Code()) }
