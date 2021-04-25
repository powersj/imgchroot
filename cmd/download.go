package cmd

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/download"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	downloadCmd = &cobra.Command{
		Use:     "download <distro> <release>",
		Example: "  imgchroot download ubuntu focal",
		Short:   "Download a cloud image",
		Args:    downloaderArgs,
		Run:     downloader,
	}
	downloadRequiredArgs int = 2
)

func init() {
	rootCmd.AddCommand(downloadCmd)
}

func downloaderArgs(cmd *cobra.Command, args []string) error {
	if len(args) != downloadRequiredArgs {
		return errors.New(
			"please provide a distro and release to download (e.g. ubuntu focal)",
		)
	}

	return nil
}

func downloader(cmd *cobra.Command, args []string) {
	imageURL, err := download.New(args[0], args[1], runtime.GOARCH)
	if err != nil {
		log.Errorln(err)
	}

	filename, err := imageURL.Download()
	if err != nil {
		log.Errorln(err)
	}

	fmt.Printf("image saved as %v \n", filename)
}
