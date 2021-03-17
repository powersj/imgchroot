package download

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUbuntu(t *testing.T) {
	download, err := New("ubuntu", "focal", "amd64")

	assert.Nil(t, err)
	assert.Contains(t, download.URL, "cloud-images.ubuntu.com")
}

func TestUnknownDistro(t *testing.T) {
	_, err := New("unknown", "focal", "amd64")
	assert.Error(t, err)
}
