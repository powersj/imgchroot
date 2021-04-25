package cmd

import (
	"errors"

	"github.com/powersj/imgchroot/pkg/configdata"
	"github.com/powersj/imgchroot/pkg/mount"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	applyCmd = &cobra.Command{
		Use:     "apply <imgage> <schema>",
		Example: "  imgchroot apply foobar.img schema.yaml",
		Short:   "Apply YAML-based configuration to the image",
		Args:    applyArgs,
		PreRun:  applyPreRun,
		Run:     apply,
	}
	applyRequiredArgs = 2
	dryRun            = false
)

func init() {
	applyCmd.Flags().IntVar(&partNum, "part-num", 1, "Partition number to mount")
	applyCmd.Flags().BoolVar(
		&dryRun, "dry-run", false, "Do not modify image and show commands that would run",
	)

	rootCmd.AddCommand(applyCmd)
}

func applyArgs(cmd *cobra.Command, args []string) error {
	if len(args) != applyRequiredArgs {
		return errors.New(
			"please provide an image and file to apply",
		)
	}

	return nil
}

func applyPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isFile(args[1])
	isRoot()
	haveRequiredCommands()
}

func apply(cmd *cobra.Command, args []string) {
	configData, err := configdata.New(args[1])
	if err != nil {
		log.Fatal(err)
	}

	configData.DryRun = dryRun
	if !configData.Validate() {
		log.Errorln("Invalid YAML")
		for _, str := range configData.SchemaErrors {
			log.Errorln(str)
		}
	}

	var image mount.Image = mount.Image{
		Filename:    args[0],
		NoSysDNS:    noSysDNS,
		NoSysMounts: noSysMounts,
		PartNum:     partNum,
		ReadOnly:    readOnly,
	}

	chroot, err := image.Mount()
	if err != nil {
		log.Errorln(err)
	}

	err = configData.Apply(chroot)
	if err != nil {
		log.Errorln(err)
	}

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
