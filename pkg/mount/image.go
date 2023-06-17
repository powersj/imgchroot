package mount

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Image is the user's main control point for manipulating the image.
type Image struct {
	Filename string `json:"filename"`

	// Settings used during mounting the image.
	PartNum     int  `json:"-"`
	NoSysDNS    bool `json:"-"`
	NoSysMounts bool `json:"-"`
	ReadOnly    bool `json:"-"`

	// Used to manipulate the image.
	MountPoint string `json:"-"`
	NBD        NBD    `json:"-"`

	// Information about the image and its' partitions.
	ImageFormat    string      `json:"image-format"`
	PartitionTable string      `json:"partition-table"`
	Size           string      `json:"size"`
	VirtualSize    string      `json:"virtual-size"`
	SectorSize     int         `json:"sector-size"`
	Partitions     []Partition `json:"partitions"`
}

// Partition captures info about each partition in an image.
type Partition struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Label      string `json:"label"`
	Filesystem string `json:"filesystem"`
}

// Info returns JSON or formatted string with details about the image.
func (i *Image) Info(jsonOutput bool) (string, error) {
	if jsonOutput {
		b, err := json.MarshalIndent(&i, "", "  ")
		if err != nil {
			return "", errors.Wrap(err, "image.Info")
		}

		return string(b), nil
	}

	var output strings.Builder = strings.Builder{}
	output.WriteString(fmt.Sprintf("%s\n", path.Base(i.Filename)))
	output.WriteString(fmt.Sprintf("type: %s\n", i.ImageFormat))
	output.WriteString(fmt.Sprintf("size: %s\n", i.Size))
	output.WriteString(fmt.Sprintf("virtual size: %s\n", i.VirtualSize))
	output.WriteString(fmt.Sprintf("partition table: %s\n", i.PartitionTable))
	output.WriteString(fmt.Sprintf("sector size: %d\n", i.SectorSize))
	output.WriteString("partitions:\n")
	for _, part := range i.Partitions {
		output.WriteString(fmt.Sprintf("  - name: %s\n", part.Name))
		output.WriteString(fmt.Sprintf("    type: %s\n", part.Type))
		output.WriteString(fmt.Sprintf("    label: %s\n", part.Label))
		output.WriteString(fmt.Sprintf("    filesystem: %s\n", part.Filesystem))
	}

	return output.String(), nil
}

// Mount does the heavy lifting to get the image ready and returns a chroot.
func (i *Image) Mount() (Chroot, error) {
	var chroot Chroot = Chroot{}

	if err := i.createMountPoint(); err != nil {
		return chroot, err
	}
	chroot.MountPoint = i.MountPoint

	if err := i.scanImage(); err != nil {
		return chroot, err
	}

	i.NBD = NBD{}
	if err := i.NBD.Connect(i.Filename, i.ImageFormat); err != nil {
		return chroot, err
	}

	if err := i.scanPartitions(); err != nil {
		return chroot, err
	}

	if err := i.mountPartition(); err != nil {
		return chroot, err
	}

	if !i.NoSysMounts {
		if err := i.mountSystemMounts(); err != nil {
			return chroot, err
		}
	}

	if !i.NoSysDNS && !i.ReadOnly {
		if err := i.mountSystemDNS(); err != nil {
			return chroot, err
		}
	}

	return chroot, nil
}

// Unmount undos all the options during mounting.
func (i *Image) Unmount() error {
	if !i.NoSysDNS && !i.ReadOnly {
		if err := i.unmountSystemDNS(); err != nil {
			return err
		}
	}

	if !i.NoSysMounts {
		if err := i.unmountSystemMounts(); err != nil {
			return err
		}
	}

	if err := exec.Command("umount", i.MountPoint).Run(); err != nil {
		return errors.Wrap(err, "image.Unmount")
	}

	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "image.Unmount")
	}

	log.Debugf("removing image mount point %s", i.MountPoint)
	os.RemoveAll(i.MountPoint)

	if err := i.NBD.Disconnect(); err != nil {
		return errors.Wrap(err, "image.Unmount")
	}

	return nil
}

// Convert bytes to MiB.
func bytes2human(bytes int64) float64 {
	return float64(bytes) / 1024 / 1024
}

// createMountPoint creates a folder in /tmp to mount the image in.
func (i *Image) createMountPoint() error {
	log.Debugf("creating temporary directory")

	dir, err := os.CreateTemp("", "imgchroot-")
	if err != nil {
		return errors.Wrap(err, "image.createMountPoint")
	}

	i.MountPoint = dir.Name()
	log.Debugf(i.MountPoint)

	return nil
}

// execLsblk runs lsblk until we get a result or timeout.
func (i *Image) execLsblk() (BlockDevice, error) {
	log.Debugln("checking for lsblk partition info")

	var device BlockDevice = BlockDevice{}
	out, err := exec.Command(
		"lsblk", "--json",
		"--output", "NAME,LOG-SEC,TYPE,PTTYPE,PARTTYPE,PARTTYPENAME,FSTYPE,LABEL",
		i.NBD.DevicePath,
	).Output()
	if err != nil {
		return device, errors.Wrap(err, "image.execLsblk")
	}

	var lsblkOutput Lsblk
	err = json.Unmarshal([]byte(out), &lsblkOutput)
	if err != nil {
		log.Errorln(err)
	}

	return lsblkOutput.LsblkDevices[0], nil
}

// execQEMUImgInfo runs qemu-img info on our filename.
func (i *Image) execQEMUImgInfo() (QEMUImgInfo, error) {
	var info QEMUImgInfo = QEMUImgInfo{}

	jsonString, err := exec.Command(
		"qemu-img", "info", "--output=json", i.Filename,
	).Output()
	if err != nil {
		return info, errors.New("qemu-img info failed")
	}

	err = json.Unmarshal([]byte(jsonString), &info)
	if err != nil {
		return info, errors.Wrap(err, "img.execQEMUImgInfo")
	}

	return info, nil
}

// mountPartition mounts and prepares the partition for use.
func (i *Image) mountPartition() error {
	var partDevice string = fmt.Sprintf("%sp%d", i.NBD.DevicePath, i.PartNum)

	var option string
	if i.ReadOnly {
		log.Debugf("mounting %s (read-only)", partDevice)
		option = "ro"
	} else {
		log.Debugf("mounting %s (read-write)", partDevice)
		option = "rw"
	}

	log.Debugf("mount -o %s %s %s", option, partDevice, i.MountPoint)
	_, err := exec.Command("mount", "-o", option, partDevice, i.MountPoint).Output()
	if err != nil {
		return errors.Wrap(err, "image.mountPartition")
	}

	log.Debugf("udevadm settle")
	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "image.mountPartition")
	}

	return nil
}

// mountSystemDNS copies the host system's DNS values in.
func (i *Image) mountSystemDNS() error {
	log.Debugf("updating resolv.conf")

	err := exec.Command(
		"mv",
		filepath.Join(i.MountPoint, "/etc/resolv.conf"),
		filepath.Join(i.MountPoint, "/etc/resolv.conf.bak"),
	).Run()
	if err != nil {
		return errors.Wrap(err, "image.mountSystemDNS")
	}

	err = exec.Command(
		"cp", "/etc/resolv.conf", filepath.Join(i.MountPoint, "/etc/resolv.conf"),
	).Run()
	if err != nil {
		return errors.Wrap(err, "image.mountSystemDNS")
	}

	return nil
}

// mountSystemMounts bind mounts /dev, /proc, and /sys.
func (i *Image) mountSystemMounts() error {
	for _, path := range []string{"/dev", "/proc", "/sys"} {
		log.Debugf("bind mounting %s", path)
		_, err := exec.Command(
			"mount", "--bind", path, filepath.Join(i.MountPoint, path),
		).Output()
		if err != nil {
			return errors.Wrap(err, "image.mountSystemMounts")
		}
	}

	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "image.mountSystemMounts")
	}

	return nil
}

// scanImage looks at the image with qemu-img info to find the format.
func (i *Image) scanImage() error {
	imgInfo, err := i.execQEMUImgInfo()
	if err != nil {
		return err
	}

	i.ImageFormat = imgInfo.Format
	i.Size = fmt.Sprintf("%.0f MiB", bytes2human(imgInfo.ActualSize))
	i.VirtualSize = fmt.Sprintf("%.0f MiB", bytes2human(imgInfo.VirtualSize))

	return nil
}

// scanPartitions looks for partitions on the nbd device.
func (i *Image) scanPartitions() error {
	blockDevice, err := i.execLsblk()
	if err != nil {
		return err
	}

	i.PartitionTable = blockDevice.Pttype
	i.SectorSize = blockDevice.LogSec

	for _, part := range blockDevice.Partitions {
		// Use the pretty name, unless empty then use the GUID value
		partType := part.Parttypename
		if strings.Compare(part.Parttypename, "") == 0 {
			partType = part.Parttype
		}

		var p Partition = Partition{
			Name:       part.Name,
			Type:       partType,
			Label:      part.Label,
			Filesystem: part.Fstype,
		}
		i.Partitions = append(i.Partitions, p)
	}

	return nil
}

// mountSystemDNS copies the host system's DNS values in.
func (i *Image) unmountSystemDNS() error {
	log.Debugf("restoring resolv.conf")

	err := exec.Command(
		"mv",
		filepath.Join(i.MountPoint, "/etc/resolv.conf.bak"),
		filepath.Join(i.MountPoint, "/etc/resolv.conf"),
	).Run()
	if err != nil {
		return errors.Wrap(err, "image.unmountSystemDNS")
	}

	return nil
}

// mountSystemMounts bind mounts /dev, /proc, and /sys.
func (i *Image) unmountSystemMounts() error {
	for _, path := range []string{"/dev", "/proc", "/sys"} {
		log.Debugf("unmounting %s", path)
		err := exec.Command("umount", filepath.Join(i.MountPoint, path)).Run()
		if err != nil {
			return errors.Wrap(err, "failed to unmount bind mounts")
		}
	}

	if err := exec.Command("udevadm", "settle").Run(); err != nil {
		return errors.Wrap(err, "image.unmountSystemMounts")
	}

	return nil
}
