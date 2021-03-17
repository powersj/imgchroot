package cmd

import (
	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/mount"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	runCmd = &cobra.Command{
		Use:     "run <image> <script>",
		Example: "  imgchroot run foobar.img setup.sh",
		Short:   "Transfer and run a file on the image",
		Args:    runArgs,
		PreRun:  runPreRun,
		Run:     run,
	}
	runRequiredArgs = 2
)

func init() {
	runCmd.Flags().BoolVar(
		&noSysDNS, "no-system-dns", false,
		"Do not drop in host system's DNS settings",
	)
	runCmd.Flags().BoolVar(
		&noSysMounts, "no-system-mounts", false,
		"Do not bind mount host system's /dev, /proc, and /sys",
	)
	runCmd.Flags().IntVar(
		&partNum, "part-num", 1,
		"Partition number to mount",
	)
	runCmd.Flags().BoolVar(
		&readOnly, "read-only", false,
		"Mount image as read-only",
	)

	rootCmd.AddCommand(runCmd)
}

func runArgs(cmd *cobra.Command, args []string) error {
	if len(args) != runRequiredArgs {
		return errors.New("please provide an image and file")
	}

	return nil
}

func runPreRun(cmd *cobra.Command, args []string) {
	isFile(args[0])
	isFile(args[1])
	isRoot()
	haveRequiredCommands()
}

func run(cmd *cobra.Command, args []string) {
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

	err = chroot.Run(args[1])
	if err != nil {
		log.Errorln(err)
	}

	err = image.Unmount()
	if err != nil {
		log.Errorln(err)
	}
}
