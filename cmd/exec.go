package cmd

import (
	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/mount"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	execCmd = &cobra.Command{
		Use:     "exec <image> [--] <command line>",
		Example: "  imgchroot exec foobar.img -- cat /etc/hosts",
		Short:   "Run a command on the image",
		Args:    execArgs,
		PreRun:  execPreRun,
		Run:     execute,
	}
	execRequiredArgs int = 2
)

func init() {
	execCmd.Flags().BoolVar(
		&noSysDNS, "no-system-dns", false,
		"Do not drop in host system's DNS settings",
	)
	execCmd.Flags().BoolVar(
		&noSysMounts, "no-system-mounts", false,
		"Do not bind mount host system's /dev, /proc, and /sys",
	)
	execCmd.Flags().IntVar(
		&partNum, "part-num", 1,
		"Partition number to mount",
	)
	execCmd.Flags().BoolVar(
		&readOnly, "read-only", false,
		"Mount image as read-only",
	)

	rootCmd.AddCommand(execCmd)
}

func execArgs(cmd *cobra.Command, args []string) error {
	if len(args) < execRequiredArgs {
		return errors.New("please provide an image and command to run")
	}

	return nil
}

func execPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isRoot()
	haveRequiredCommands()
}

func execute(cmd *cobra.Command, args []string) {
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

	err = chroot.Exec(args[1:])
	if err != nil {
		log.Errorln(err)
	}

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
