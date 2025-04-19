package storage

import "github.com/heathcliff26/go-wol/pkg/server/storage/file"

const (
	// TODO: Make default true when implementing api
	DEFAULT_READONLY     = true
	DEFAULT_BACKEND_TYPE = "file"
)

type StorageConfig struct {
	Type     string                 `json:"type"`
	Readonly bool                   `json:"-,omitempty"`
	File     file.FileBackendConfig `json:"file,omitempty"`
}

func NewDefaultStorageConfig() StorageConfig {
	return StorageConfig{
		Type:     DEFAULT_BACKEND_TYPE,
		Readonly: DEFAULT_READONLY,
		File:     file.NewDefaultFileBackendConfig(),
	}
}
