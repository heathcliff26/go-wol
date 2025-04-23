package storage

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"sync"

	"github.com/heathcliff26/go-wol/pkg/server/storage/file"
	"github.com/heathcliff26/go-wol/pkg/server/storage/types"
	"github.com/heathcliff26/go-wol/pkg/server/storage/valkey"
	"github.com/heathcliff26/go-wol/pkg/version"
	"github.com/heathcliff26/go-wol/static"
)

type Storage struct {
	backend  types.StorageBackend
	readonly bool

	indexLock     sync.RWMutex
	indexHTML     string
	indexChecksum string
}

func NewStorage(cfg StorageConfig) (*Storage, error) {
	var backend types.StorageBackend
	var err error
	switch cfg.Type {
	case "file":
		backend, err = file.NewFileBackend(cfg.File)
	case "valkey":
		backend, err = valkey.NewValkeyBackend(cfg.Valkey)
	default:
		return nil, fmt.Errorf("unknown storage backend type: %s", cfg.Type)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create storage backend: %w", err)
	}

	s := &Storage{
		backend:  backend,
		readonly: cfg.Readonly,
	}

	if !s.readonly {
		s.readonly, err = s.backend.Readonly()
		if err != nil {
			return nil, fmt.Errorf("failed to check if storage backend is readonly: %w", err)
		}
	}

	err = s.updateIndexHTML()
	if err != nil {
		return nil, fmt.Errorf("failed to create index.html: %w", err)
	}

	return s, nil
}

type indexValues struct {
	Readonly bool
	Hosts    []types.Host
	Version  string
	Name     string
}

// Generate the index.html file from the template and the current hosts
func (s *Storage) updateIndexHTML() error {
	s.indexLock.Lock()
	defer s.indexLock.Unlock()

	hosts, err := s.backend.GetHosts()
	if err != nil {
		return fmt.Errorf("failed to get hosts: %w", err)
	}

	values := indexValues{
		Readonly: s.readonly,
		Hosts:    hosts,
		Version:  version.Version(),
		Name:     version.Name,
	}

	tmpl, err := template.New("index.html").Parse(string(static.IndexTemplate))
	if err != nil {
		return fmt.Errorf("unexpected error when creating a template from static html: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, values)
	if err != nil {
		return fmt.Errorf("failed to apply values to html template: %w", err)
	}
	indexHTML := buf.String()

	checksum := sha256.Sum256([]byte(indexHTML))

	s.indexHTML = indexHTML
	s.indexChecksum = hex.EncodeToString(checksum[:])

	return nil
}

// Return the current index.html and its checksum
func (s *Storage) GetIndexHTML() (string, string) {
	s.indexLock.RLock()
	defer s.indexLock.RUnlock()

	return s.indexHTML, s.indexChecksum
}

// Return if the storage is readonly
func (s *Storage) Readonly() bool {
	return s.readonly
}

// Add a new host and update the index.html
func (s *Storage) AddHost(mac, host string) error {
	if s.readonly {
		return fmt.Errorf("storage is readonly")
	}

	err := s.backend.AddHost(mac, host)
	if err != nil {
		return fmt.Errorf("failed to add host: %w", err)
	}

	err = s.updateIndexHTML()
	if err != nil {
		return err
	}

	return nil
}

// Remove a host and update the index.html
func (s *Storage) RemoveHost(mac string) error {
	if s.readonly {
		return fmt.Errorf("storage is readonly")
	}

	err := s.backend.RemoveHost(mac)
	if err != nil {
		return fmt.Errorf("failed to remove host: %w", err)
	}

	err = s.updateIndexHTML()
	if err != nil {
		return err
	}

	return nil
}
