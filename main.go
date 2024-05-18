package main

import (
	"github.com/iliyankg/colab-shield/cli/cmd"
	"github.com/iliyankg/colab-shield/cli/config"
)

func main() {
	config.Init()
	cmd.Execute()
}
