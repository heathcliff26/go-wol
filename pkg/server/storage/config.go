package storage

import (
	"github.com/heathcliff26/go-wol/pkg/server/storage/file"
	"github.com/heathcliff26/go-wol/pkg/server/storage/valkey"
)

const (
	DEFAULT_READONLY     = false
	DEFAULT_BACKEND_TYPE = "file"
)

type StorageConfig struct {
	Type     string                 `json:"type"`
	Readonly bool                   `json:"readonly,omitempty"`
	File     file.FileBackendConfig `json:"file,omitempty"`
	Valkey   valkey.ValkeyConfig    `json:"valkey,omitempty"`
}

func NewDefaultStorageConfig() StorageConfig {
	return StorageConfig{
		Type:     DEFAULT_BACKEND_TYPE,
		Readonly: DEFAULT_READONLY,
		File:     file.NewDefaultFileBackendConfig(),
	}
}
