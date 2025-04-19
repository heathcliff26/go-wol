package file

import "os"

const (
	DEFAULT_FILE_PATH           = "hosts.yaml"
	DEFAULT_FILE_PATH_CONTAINER = "/data/hosts.yaml"
)

type FileBackendConfig struct {
	Path string `json:"path,omitempty"`
}

func NewDefaultFileBackendConfig() FileBackendConfig {
	cfg := FileBackendConfig{
		Path: DEFAULT_FILE_PATH,
	}
	if _, ok := os.LookupEnv("container"); ok {
		cfg.Path = DEFAULT_FILE_PATH_CONTAINER
	}
	return cfg
}
