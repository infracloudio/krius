package spec

import (
	"fmt"
	"log"

	spec "github.com/infracloudio/krius/pkg/specvalidate"
	"github.com/spf13/cobra"
)

const (
	configFile = "config-file"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Applies/Updates the give profie file",
	Run:   applySpec,
}

func init() {
	specCmd.AddCommand(applyCmd)
	addSpecApplyFlags(applyCmd)
}

func applySpec(cmd *cobra.Command, args []string) {
	configFileFlag, _ := cmd.Flags().GetString(configFile)
	loader, ruleSchemaLoader, err := spec.GetLoaders(configFileFlag)
	if err != nil {
		log.Println(err)
		return
	}
	valid, errors := spec.ValidateYML(loader, ruleSchemaLoader)
	if !valid {
		log.Println(errors)
		return
	}
	fmt.Println("valid yaml")
}
