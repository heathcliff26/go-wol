package storage

import (
	"os"
	"testing"

	"github.com/heathcliff26/go-wol/pkg/server/storage/file"
	"github.com/heathcliff26/go-wol/pkg/server/storage/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBackend struct {
	mock.Mock
	types.StorageBackend
}

func (m *MockBackend) AddHost(mac, host string) error {
	args := m.Called(mac, host)
	return args.Error(0)
}

func (m *MockBackend) RemoveHost(mac string) error {
	args := m.Called(mac)
	return args.Error(0)
}

func (m *MockBackend) GetHost(mac string) (string, error) {
	args := m.Called(mac)
	return args.String(0), args.Error(1)
}

func (m *MockBackend) GetHosts() ([]types.Host, error) {
	args := m.Called()
	return args.Get(0).([]types.Host), args.Error(1)
}

func (m *MockBackend) Readonly() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func TestNewStorage(t *testing.T) {
	t.Run("FileBackend", func(t *testing.T) {
		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			Readonly: true,
		}

		s, err := NewStorage(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, s)
		assert.IsType(t, &file.FileBackend{}, s.backend, "Should have file as backend type")
		assert.Equal(t, cfg.Readonly, s.readonly, "Readonly should match config")
		assert.NotEmpty(t, s.indexHTML, "Index HTML should not be empty")
		assert.NotEmpty(t, s.indexChecksum, "Index checksum should not be empty")
	})
	t.Run("UnknownBackend", func(t *testing.T) {
		s, err := NewStorage(StorageConfig{Type: "unknown"})
		assert.Nil(t, s, "Storage should be nil")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown storage backend type")
	})

	t.Run("WritableBackend", func(t *testing.T) {
		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: t.TempDir() + "/test.yaml",
			},
			Readonly: false,
		}

		s, err := NewStorage(cfg)
		assert.NoError(t, err, "Should not return an error")
		assert.NotNil(t, s, "Storage should not be nil")

		assert.False(t, s.Readonly())
	})

	t.Run("OverwriteConfigWhenReadonlyBackend", func(t *testing.T) {
		path := t.TempDir() + "/test.yaml"

		f, err := os.OpenFile(path, os.O_CREATE|os.O_RDONLY, 0444)
		if !assert.NoError(t, err, "Should create file") {
			t.FailNow()
		}
		f.Close()

		cfg := StorageConfig{
			Type: "file",
			File: file.FileBackendConfig{
				Path: path,
			},
			Readonly: false,
		}

		s, err := NewStorage(cfg)
		assert.NoError(t, err, "Should not return an error")
		assert.NotNil(t, s, "Storage should not be nil")

		assert.True(t, s.Readonly(), "Should ignore config if backend is readonly")
	})
}

func TestStorageGetIndexHTML(t *testing.T) {
	s := &Storage{
		indexHTML:     "<html>Test</html>",
		indexChecksum: "1234567890abcdef",
	}

	html, checksum := s.GetIndexHTML()
	assert.Equal(t, s.indexHTML, html, "HTML should match")
	assert.Equal(t, s.indexChecksum, checksum, "Checksum should match")
}

func TestStorageReadonly(t *testing.T) {
	s := &Storage{
		readonly: true,
	}

	assert.True(t, s.Readonly())

	s.readonly = false
	assert.False(t, s.Readonly())
}

func TestStorageAddHost(t *testing.T) {
	mockBackend := new(MockBackend)
	s := &Storage{
		backend:  mockBackend,
		readonly: false,
	}

	t.Run("Readonly", func(t *testing.T) {
		s.readonly = true
		err := s.AddHost("00:11:22:33:44:55", "test")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage is readonly")
	})

	t.Run("Success", func(t *testing.T) {
		assert := assert.New(t)

		s.readonly = false
		mockBackend.On("AddHost", "00:11:22:33:44:55", "test").Return(nil)
		mockBackend.On("GetHosts").Return([]types.Host{{MAC: "00:11:22:33:44:55", Name: "test"}}, nil)

		err := s.AddHost("00:11:22:33:44:55", "test")
		assert.NoError(err, "Should add host without error")
		assert.NotEmpty(s.indexHTML, "Should have updated index HTML")
		assert.NotEmpty(s.indexChecksum, "Should have updated index checksum")
		mockBackend.AssertExpectations(t)
	})
}

func TestStorageRemoveHost(t *testing.T) {
	mockBackend := new(MockBackend)
	s := &Storage{
		backend:  mockBackend,
		readonly: false,
	}

	t.Run("ReadonlyStorage", func(t *testing.T) {
		s.readonly = true
		err := s.RemoveHost("00:11:22:33:44:55")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "storage is readonly")
	})

	t.Run("SuccessfulRemove", func(t *testing.T) {
		assert := assert.New(t)

		s.readonly = false
		mockBackend.On("RemoveHost", "00:11:22:33:44:55").Return(nil)
		mockBackend.On("GetHosts").Return([]types.Host{}, nil)

		err := s.RemoveHost("00:11:22:33:44:55")
		assert.NoError(err, "Should remove host without error")
		assert.NotEmpty(s.indexHTML, "Should have updated index HTML")
		assert.NotEmpty(s.indexChecksum, "Should have updated index checksum")
		mockBackend.AssertExpectations(t)
	})
}
