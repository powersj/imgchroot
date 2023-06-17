package download

import (
	"fmt"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	bytesToMB          int64         = 1024 * 1024
	progressPercent    float64       = 100
	progressPrintDelay time.Duration = 500
)

// ImageURL captures the details and URL of the image to download.
type ImageURL struct {
	Distro  string
	Release string
	Arch    string
	URL     string
}

// New initializes a new ImageURL struct to create a download URL.
func New(distro string, release string, arch string) (*ImageURL, error) {
	imageURL := &ImageURL{
		Distro:  distro,
		Release: release,
		Arch:    arch,
	}

	if err := imageURL.build(); err != nil {
		return nil, errors.Wrap(err, "unable to create download URL")
	}

	return imageURL, nil
}

// Download an image.
func (i *ImageURL) Download() (string, error) {
	fmt.Printf("downloading %v\n", i.URL)
	client := grab.NewClient()
	req, _ := grab.NewRequest(".", i.URL)

	resp := client.Do(req)
	log.Debugf("  %v\n", resp.HTTPResponse.Status)

	ticker := time.NewTicker(progressPrintDelay * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			printStatus(resp)
		case <-resp.Done:
			printStatus(resp)
			fmt.Println("")

			if err := resp.Err(); err != nil {
				return "", errors.Wrap(err, "download failed")
			}

			return resp.Filename, nil
		}
	}
}

// build a URL to download from.
func (i *ImageURL) build() error {
	if i.Distro == "ubuntu" {
		return i.ubuntuURL()
	}

	return errors.New(fmt.Sprintf("unknown distro: %s", i.Distro))
}

// ubuntuURL stores details about Ubuntu URLs.
func (i *ImageURL) ubuntuURL() error {
	var hostname string = "https://cloud-images.ubuntu.com"

	i.URL = fmt.Sprintf(
		"%s/%s/current/%s-server-cloudimg-%s.img", hostname, i.Release, i.Release, i.Arch,
	)

	return nil
}

// printStatus prints the current download status.
//
//	Example: " 532 of  532 MB (100.0%) @  10.3 MB/sec"
func printStatus(resp *grab.Response) {
	fmt.Printf(
		"\r%4v of %4v MB (%4.1f%%) @ %5.1f MB/sec",
		resp.BytesComplete()/bytesToMB,
		resp.Size()/bytesToMB,
		progressPercent*resp.Progress(),
		resp.BytesPerSecond()/float64(bytesToMB),
	)
}
