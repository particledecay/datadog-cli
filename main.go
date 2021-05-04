package main

import (
	"github.com/particledecay/datadog-cli/cmd"
	"github.com/particledecay/datadog-cli/config"
)

func main() {
	config.Load()
	cmd.Execute()
}
