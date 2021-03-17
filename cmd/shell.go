package cmd

import (
	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/mount"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	shellCmd = &cobra.Command{
		Use:     "shell <image>",
		Example: "  imgchroot shell foobar.img",
		Short:   "Start a shell on the image",
		Args:    shellArgs,
		PreRun:  shellPreRun,
		Run:     shellRun,
	}
	shell string = "/bin/bash"
)

func init() {
	shellCmd.Flags().BoolVar(
		&noSysDNS, "no-system-dns", false,
		"Do not drop in host system's DNS settings",
	)
	shellCmd.Flags().BoolVar(
		&noSysMounts, "no-system-mounts", false,
		"Do not bind mount host system's /dev, /proc, and /sys",
	)
	shellCmd.Flags().IntVar(
		&partNum, "part-num", 1,
		"Partition number to mount",
	)
	shellCmd.Flags().BoolVar(
		&readOnly, "read-only", false,
		"Mount image as read-only",
	)

	rootCmd.AddCommand(shellCmd)
}

func shellArgs(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("please provide an image to operate on")
	}

	return nil
}

func shellPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isRoot()
	haveRequiredCommands()
}

func shellRun(cmd *cobra.Command, args []string) {
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

	err = chroot.Shell(shell)
	if err != nil {
		log.Errorln(err)
	}

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
