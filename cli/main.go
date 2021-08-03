package main

import (
	"github.com/infracloudio/krius/cli/cmd"
	_ "github.com/infracloudio/krius/cli/cmd/spec"
	randomSeed "github.com/infracloudio/krius/pkg/utils"
)

func main() {
	randomSeed.Init()
	cmd.Execute()
}
