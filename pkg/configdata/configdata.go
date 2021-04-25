package configdata

import (
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/pkg/errors"
	"github.com/powersj/imgchroot/pkg/mount"
	"gopkg.in/validator.v2"
	"gopkg.in/yaml.v3"
)

// ConfigData captures the config data passed in and any errors.
type ConfigData struct {
	DryRun       bool
	Filename     string
	Schema       Schema
	SchemaErrors []string
}

// New initializes a new ConfigData struct.
func New(filename string) (*ConfigData, error) {
	configData := &ConfigData{
		Filename: filename,
	}

	if err := configData.parseFile(configData.Filename); err != nil {
		return nil, errors.Wrap(err, "parsing file failed")
	}

	return configData, nil
}

// Apply config data against a given chroot.
func (c *ConfigData) Apply(chroot mount.Chroot) error {
	// TODO:
	// determine OS type
	// case on OS type create object with functions to do each thing

	// TODO:
	// keeps list of errors in each phase
	// apply returns all the errors?

	return nil
}

// Validate a passed in config data file.
func (c *ConfigData) Validate() bool {
	err := validator.Validate(c.Schema)
	if err == nil {
		return true
	}

	// capture all errors in a sorted array
	for field, errorStr := range err.(validator.ErrorMap) {
		c.SchemaErrors = append(c.SchemaErrors, fmt.Sprintf("%s - %s", field, errorStr))
	}
	sort.Strings(c.SchemaErrors)

	return false
}

// parseFile will attempt to parse the file to the schema.
func (c *ConfigData) parseFile(configDataFile string) error {
	rawConfigData, err := ioutil.ReadFile(configDataFile)
	if err != nil {
		return errors.Wrap(err, "unable to read file")
	}

	c.Schema = Schema{}
	err = yaml.Unmarshal(rawConfigData, &c.Schema)
	if err != nil {
		return errors.Wrap(err, "unable to unmarshal YAML")
	}

	return nil
}
