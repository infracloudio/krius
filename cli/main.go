package main

import (
	"github.com/infracloudio/krius/cli/cmd"
	_ "github.com/infracloudio/krius/cli/cmd/spec"
	random "github.com/infracloudio/krius/pkg/random"
)

func main() {
	random.Init()

	cmd.Execute()
}
