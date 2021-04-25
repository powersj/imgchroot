package configdata

// Schema captures all configuration data options.
type Schema struct {
	Version       int               `yaml:"version" validate:"nonzero"`
	Variables     map[string]string `yaml:"variables,omitempty"`
	EarlyCommands [][]string        `yaml:"early-commands,omitempty"`
	Archives      []ArchivesStruct  `yaml:"sources,omitempty"`
	PPA           []PPAStruct       `yaml:"ppa,omitempty"`
	Packages      struct {
		Install  []string `yaml:"install,omitempty"`
		Remove   []string `yaml:"remove,omitempty"`
		Proposed []string `yaml:"proposed,omitempty"`
	} `yaml:"packages,omitempty"`
	Snap     []SnapStruct  `yaml:"snap,omitempty"`
	Files    []FilesStruct `yaml:"files,omitempty"`
	Groups   []string      `yaml:"groups,omitempty"`
	Users    []UsersStruct `yaml:"users,omitempty"`
	Commands [][]string    `yaml:"commands,omitempty"`
}

// ArchivesStruct captures data for adding archives.
type ArchivesStruct struct {
	Name    string `yaml:"name" validate:"nonzero"`
	URL     string `yaml:"url" validate:"nonzero"`
	Comment string `yaml:"comment,omitempty"`
	Remove  bool   `yaml:"remove,omitempty"`
}

// FilesStruct captures data for creating files.
type FilesStruct struct {
	Content     string `yaml:"content" validate:"nonzero"`
	Path        string `yaml:"path" validate:"nonzero"`
	Permissions string `yaml:"permissions,omitempty"`
	Owner       string `yaml:"owner,omitempty"`
}

// PPAStruct captures data for installing a PPA.
type PPAStruct struct {
	URL         string `yaml:"url" validate:"nonzero"`
	Fingerprint string `yaml:"fingerprint,omitempty"`
	Creds       string `yaml:"creds,omitempty"`
	Comment     string `yaml:"comment,omitempty"`
	Remove      bool   `yaml:"remove,omitempty"`
}

// SnapStruct captures data for installing any snaps.
type SnapStruct struct {
	Name    string `yaml:"name" validate:"nonzero"`
	Channel string `yaml:"channel" validate:"nonzero"`
}

// UsersStruct captures data for new users.
type UsersStruct struct {
	Username          string   `yaml:"username" validate:"nonzero"`
	Password          string   `yaml:"password,omitempty"`
	SSHAuthorizedKeys []string `yaml:"ssh-authorized-keys,omitempty"`
	SSHImportID       []string `yaml:"ssh-import-id,omitempty"`
	Groups            []string `yaml:"groups,omitempty"`
	Sudoer            bool     `yaml:"sudoer,omitempty"`
}
