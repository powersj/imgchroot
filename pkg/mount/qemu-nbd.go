package mount

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// NBD collects details about the used NBD device.
type NBD struct {
	// /sys/block/nbd0
	BlockPath string
	// /dev/nbd0
	DevicePath string
	// nbd0
	Name string
	// PID of the nbd device
	PID string
	// /sys/block/nbd0/pid
	PIDFile string
}

// Connect a given image.
func (n *NBD) Connect(image string, format string) error {
	if err := n.allocate(); err != nil {
		return err
	}

	var cmd *exec.Cmd = exec.Command(
		"qemu-nbd",
		fmt.Sprintf("--format=%s", format),
		"--connect", n.DevicePath,
		image,
	)

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "nbd.Connect")
	}

	if err := n.waitForPID(); err != nil {
		return errors.Wrap(err, "nbd.Connect")
	}

	log.Debugln("blockdev re-read partition table")
	if err := exec.Command("blockdev", "--rereadpt", n.DevicePath).Run(); err != nil {
		return errors.Wrap(err, "nbd.Connect")
	}

	log.Debugln("udevadm settle")
	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "nbd.Connect")
	}

	return nil
}

// Disconnect the NBD device from qemu-nbd to free it.
func (n *NBD) Disconnect() error {
	log.Debugf("disconnecting image from %s", n.Name)

	if err := exec.Command("qemu-nbd", "--disconnect", n.DevicePath).Run(); err != nil {
		return errors.Wrap(err, "nbd.Disconnect")
	}

	if err := n.waitForPIDCleanup(); err != nil {
		return errors.Wrap(err, "nbd.Disconnect")
	}

	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "nbd.Disconnect")
	}

	return nil
}

// Allocate checks for nbd kernel module and finds an empty NBD device to use.
func (n *NBD) allocate() error {
	log.Debugf("allocating nbd device")

	if !n.isNBDLoaded() {
		if err := n.loadNBD(); err != nil {
			return errors.Wrap(err, "nbd.allocate")
		}
	}

	files, err := filepath.Glob("/sys/block/nbd*")
	if err != nil {
		return errors.Wrap(err, "nbd.allocate")
	}

	for _, file := range files {
		if _, err := os.Stat(path.Join(file, "pid")); os.IsNotExist(err) {
			n.Name = filepath.Base(file)
			n.DevicePath = path.Join("/dev", filepath.Base(file))
			n.BlockPath = file
			n.PIDFile = filepath.Join(file, "pid")

			log.Debugf(n.Name)

			return nil
		}
	}

	return errors.New("Unable to allocate an NBD device")
}

// isNBDLoaded verifies if the nbd kernel module is loaded.
func (n *NBD) isNBDLoaded() bool {
	out, err := exec.Command("lsmod").Output()
	if err != nil {
		log.Errorln(err)
	}

	if !strings.Contains(string(out), "nbd") {
		return false
	}

	return true
}

// loadNBD loads the NBD kernel module if required.
func (n *NBD) loadNBD() error {
	log.Debugf("loading nbd kernel module")

	if err := exec.Command("modprobe", "nbd").Run(); err != nil {
		return errors.Wrap(err, "nbd.loadNBD")
	}

	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "nbd.loadNBD")
	}

	log.Debugf("nbd kernel module loaded")

	return nil
}

// waitForPID waits on pidfile for nbd device.
func (n *NBD) waitForPID() error {
	for i := 0; i < 30; i++ {
		log.Debugln("checking for pid file")

		data, err := os.ReadFile(n.PIDFile)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		n.PID = strings.TrimSpace(string(data))
		log.Debugf(n.PID)
		return nil
	}

	return errors.New("timed out waiting for nbd file")
}

// waitForPIDCleanup need to wait for the PID file to disappear.
func (n *NBD) waitForPIDCleanup() error {
	for i := 0; i < 30; i++ {
		log.Debugln("checking for pid file")

		_, err := os.ReadFile(n.PIDFile)
		if err != nil {
			return nil
		}

		time.Sleep(time.Second)
		continue
	}

	return errors.New("timed out waiting for nbd file to cleanup")
}
