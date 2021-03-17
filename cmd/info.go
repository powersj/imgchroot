package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/mount"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:     "info <image>",
	Example: "  imgchroot info foobar.img",
	Short:   "Information about the image",
	Args:    infoArgs,
	PreRun:  infoPreRun,
	Run:     info,
}

func init() {
	infoCmd.Flags().BoolVar(&jsonOutput, "json", false, "Print output in JSON")

	rootCmd.AddCommand(infoCmd)
}

func infoArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("please provide an image to operate on")
	}

	return nil
}

func infoPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isRoot()
	haveRequiredCommands()
}

func info(cmd *cobra.Command, args []string) {
	var image mount.Image = mount.Image{
		Filename:    args[0],
		NoSysDNS:    true,
		NoSysMounts: true,
		PartNum:     partNum,
		ReadOnly:    true,
	}

	_, err := image.Mount()
	if err != nil {
		log.Errorln(err)
	}

	output, err := image.Info(jsonOutput)
	if err != nil {
		log.Errorln(err)
	}

	fmt.Print(output)

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
