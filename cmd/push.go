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
	pushCmd = &cobra.Command{
		Use:     "push <image> <source path> <destination path>",
		Example: "  imgchroot push foobar.img custom.list /etc/apt/sources.list.d/",
		Short:   "Push a file to the image",
		Args:    pushArgs,
		PreRun:  pushPreRun,
		Run:     push,
	}
	pushRequiredArgs int = 3
)

func init() {
	pushCmd.Flags().IntVar(&partNum, "part-num", 1, "Partition number to mount")

	rootCmd.AddCommand(pushCmd)
}

// Checks that there is exactly one arguments.
func pushArgs(cmd *cobra.Command, args []string) error {
	if len(args) != pushRequiredArgs {
		return errors.New(
			"please provide an image to operate on, a source, and destination",
		)
	}

	if _, err := os.Stat(args[2]); err != nil {
		return errors.New("source does not exist")
	}

	if !filepath.IsAbs(args[2]) {
		return errors.New("absolute path required for destination")
	}

	return nil
}

func pushPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isRoot()
	haveRequiredCommands()
}

func push(cmd *cobra.Command, args []string) {
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

	err = chroot.Push(args[1], args[2])
	if err != nil {
		log.Errorln(err)
	}

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
