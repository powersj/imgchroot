package cmd

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/mount"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	pullCmd = &cobra.Command{
		Use:     "pull <image> <source path> <destination path>",
		Example: "  imgchroot pull foobar.img /etc/hosts .",
		Short:   "Pull a file from the image",
		Args:    pullArgs,
		PreRun:  pullPreRun,
		Run:     pull,
	}
	pullRequiredArgs int = 3
)

func init() {
	pullCmd.Flags().IntVar(&partNum, "part-num", 1, "Partition number to mount")

	rootCmd.AddCommand(pullCmd)
}

func pullArgs(cmd *cobra.Command, args []string) error {
	if len(args) != pullRequiredArgs {
		return errors.New(
			"please provide an image to operate on, a source, and destination",
		)
	}

	if !filepath.IsAbs(args[1]) {
		return errors.New("absolute path required for source")
	}

	if _, err := os.Stat(args[2]); err == nil {
		return errors.New("destination already exists")
	}

	return nil
}

func pullPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isRoot()
	haveRequiredCommands()
}

func pull(cmd *cobra.Command, args []string) {
	var image mount.Image = mount.Image{
		Filename:    args[0],
		NoSysDNS:    true,
		NoSysMounts: true,
		PartNum:     partNum,
		ReadOnly:    readOnly,
	}

	chroot, err := image.Mount()
	if err != nil {
		log.Errorln(err)
	}

	err = chroot.Pull(args[1], args[2])
	if err != nil {
		log.Errorln(err)
	}

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
