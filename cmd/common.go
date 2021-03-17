package cmd

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

const (
	version = "v1.1.0"
)

var (
	debugOutput      bool     = false
	jsonOutput       bool     = false
	noSysDNS         bool     = false
	noSysMounts      bool     = false
	partNum          int      = 1
	readOnly         bool     = false
	requiredCommands []string = []string{"qemu-img", "qemu-nbd"}
)

// Check for required commands.
func haveRequiredCommands() {
	log.Debugf("checking for required commands")
	for _, command := range requiredCommands {
		if err := exec.Command("which", command).Run(); err != nil {
			log.Fatalf("please install %s to run", command)
		}
	}
}

// Verify filename is a file that exists and not a directory.
func isFile(filename string) {
	stat, err := os.Stat(filename)

	if os.IsNotExist(err) {
		log.Fatalf("%s: does not exist", filename)
	} else if stat.IsDir() {
		log.Fatalf("%s: is a directory", filename)
	}
}

// Check if running as root.
func isRoot() {
	if os.Geteuid() != 0 {
		log.Fatal("root privilege required to run")
	}
}
