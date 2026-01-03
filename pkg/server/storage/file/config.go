package file

const (
	DEFAULT_FILE_PATH = "hosts.yaml"
)

type FileBackendConfig struct {
	Path string `yaml:"path,omitempty"`
}

func NewDefaultFileBackendConfig() FileBackendConfig {
	return FileBackendConfig{
		Path: DEFAULT_FILE_PATH,
	}
}
