package main

import (
	"fmt"

	"github.com/infracloudio/krius/cli/cmd"
	_ "github.com/infracloudio/krius/cli/cmd/spec"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./..")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	cmd.Execute()
}
