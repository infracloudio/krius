package main

import (
	"github.com/infracloudio/krius/cli/cmd"
	_ "github.com/infracloudio/krius/cli/cmd/spec"
)

func main() {
	cmd.Execute()
}
