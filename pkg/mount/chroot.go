package mount

import (
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Chroot details about the chroot storing the image.
type Chroot struct {
	MountPoint string
}

// Exec a command in the image.
func (c *Chroot) Exec(args []string) error {
	args = append([]string{c.MountPoint}, args...)
	cmd := exec.Command("chroot", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_ = cmd.Run()

	return nil
}

// Pull takes a file from the image to the users' system.
func (c *Chroot) Pull(srcPath string, destPath string) error {
	err := exec.Command("cp", filepath.Join(c.MountPoint, srcPath), destPath).Run()
	if err != nil {
		log.Errorln(err)
	}

	return nil
}

// Push takes a file from the user's system to the image.
func (c *Chroot) Push(srcPath string, destPath string) error {
	err := exec.Command("cp", srcPath, filepath.Join(c.MountPoint, destPath)).Run()
	if err != nil {
		log.Errorln(err)
	}

	return nil
}

// Run a command on the image.
func (c *Chroot) Run(scriptPath string) error {
	err := c.Push(scriptPath, "/tmp/")
	if err != nil {
		log.Errorln(err)
	}

	remoteScript := filepath.Join("/tmp/", filepath.Base(scriptPath))
	err = c.Exec([]string{"chmod", "755", remoteScript})
	if err != nil {
		log.Errorln(err)
	}
	err = c.Exec([]string{remoteScript})
	if err != nil {
		log.Errorln(err)
	}

	err = c.Exec([]string{"rm", remoteScript})
	if err != nil {
		log.Errorln(err)
	}

	return nil
}

// Shell into the image.
func (c *Chroot) Shell(shell string) error {
	err := c.Exec([]string{shell})
	if err != nil {
		log.Errorln(err)
	}

	return nil
}
