package mount

// QEMUImgInfo captures qemu-img JSON output.
//
//	e.g. qemu-img info --output=json focal-server-cloudimg-amd64.img
type QEMUImgInfo struct {
	// Filename analyzed
	Filename string `json:"filename"`

	// Type of image (e.g. qcow2, raw, etc.)
	Format string `json:"format"`

	// Physical size of the disk (e.g. 550637568)
	ActualSize int64 `json:"actual-size"`

	// Virtual size of the disk (e.g. 2361393152)
	VirtualSize int64 `json:"virtual-size"`
}
