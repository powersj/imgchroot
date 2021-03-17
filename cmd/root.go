package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "imgchroot <command>",
	Short: "Quickly interact and customize a cloud image",
	Long: `Quickly interact and customize a cloud image

imgchroot is a Go-based CLI to quickly customize cloud images in
a chroot without the need to boot the image or setup a user with
credentials.

imgchroot mounts the image to a temporary directory, using the
network block device (NBD) protocol. It then runs the required
operation, such as a command, moving a file, or starts a shell
in the chroot. Once the operation is complete the image is
unmounted all without needing to boot the image itself.`,
	PersistentPreRunE: setup,
}

// CLI function to setup flags.
func init() {
	rootCmd.Version = version

	rootCmd.PersistentFlags().BoolVar(
		&debugOutput, "debug", false, "debug output",
	)
}

// Called before all commands to setup general run-time settings.
func setup(cmd *cobra.Command, args []string) error {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.999999999",
	})

	if debugOutput {
		log.SetLevel(log.DebugLevel)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags.
//
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
