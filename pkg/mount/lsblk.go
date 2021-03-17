package mount

// Lsblk captures the JSON output of lsblk command.
type Lsblk struct {
	LsblkDevices []BlockDevice `json:"blockdevices"`
}

// BlockDevice a unique device in lsblk output.
type BlockDevice struct {
	// nbd1
	Name string `json:"name"`
	// 512
	LogSec int `json:"log-sec"`
	// gpt
	Pttype string `json:"pttype"`
	// Child partitions of a device.
	Partitions []Child `json:"children"`
}

// Child details each partition of a device.
type Child struct {
	// nbd1p1
	Name string `json:"name"`
	// 0fc63daf-8483-4772-8e79-3d69d8477de4
	Parttype string `json:"parttype"`
	// Linux filesystem
	Parttypename string `json:"parttypename"`
	// ext4
	Fstype string `json:"fstype"`
	// cloudimg-rootfs
	Label string `json:"label"`
}
