package cmd

import (
	"errors"

	"github.com/powersj/imgchroot/pkg/configdata"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:     "validate <schema>",
	Example: "  imgchroot validate schema.yaml",
	Short:   "Validate YAML-based configuration files",
	Args:    validateArgs,
	PreRun:  validatePreRun,
	Run:     validate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New(
			"please provide a file to validate",
		)
	}

	return nil
}

func validatePreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
}

func validate(cmd *cobra.Command, args []string) {
	configData, err := configdata.New(args[0])
	if err != nil {
		log.Fatal(err)
		return
	}

	if !configData.Validate() {
		for _, str := range configData.SchemaErrors {
			log.Errorln(str)
		}
	}
}
