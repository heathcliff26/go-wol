package testsuite

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Runs basic tests for the storage backend.
// It will create a new backend instance for each test case.
func RunStorageBackendTests(t *testing.T, factory StorageBackendFactory) {
	t.Run("AddHost", func(t *testing.T) {
		backend := factory(t, "add-host")

		err := backend.AddHost(testHosts[0].MAC, testHosts[0].Name)
		assert.NoError(t, err, "AddHost failed")
	})

	t.Run("GetHost", func(t *testing.T) {
		backend := factory(t, "get-host")
		addHosts(t, backend)

		host, err := backend.GetHost(testHosts[0].MAC)
		assert.NoError(t, err, "GetHost failed")
		assert.Equal(t, testHosts[0].Name, host, "Failed to retrieve host")
	})

	t.Run("RemoveHost", func(t *testing.T) {
		tMatrix := []struct {
			name  string
			index int
		}{
			{
				name:  "FirstElement",
				index: 0,
			},
			{
				name:  "LastElement",
				index: len(testHosts) - 1,
			},
			{
				name:  "MiddleElement",
				index: 1,
			},
		}

		for i, tCase := range tMatrix {
			t.Run(tCase.name, func(t *testing.T) {
				backend := factory(t, fmt.Sprintf("remove-host-%d", i))
				addHosts(t, backend)

				err := backend.RemoveHost(testHosts[tCase.index].MAC)
				assert.NoError(t, err, "RemoveHost failed")
				host, _ := backend.GetHost(testHosts[tCase.index].MAC)
				assert.Empty(t, host, "Expected empty host")
			})
		}
	})

	t.Run("RemoveHostNonExistent", func(t *testing.T) {
		backend := factory(t, "remove-host-non-existent")

		err := backend.RemoveHost("00:11:22:33:44:55")
		assert.NoError(t, err, "RemoveHost for non-existent host failed")
	})

	t.Run("GetHosts", func(t *testing.T) {
		backend := factory(t, "get-hosts")
		addHosts(t, backend)

		hosts, err := backend.GetHosts()
		assert.NoError(t, err, "GetHosts failed")
		assert.Equal(t, testHosts, hosts, "Expected result to match testHosts")
	})
}
