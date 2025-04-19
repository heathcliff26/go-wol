//go:build !race

package testsuite

import (
	"testing"
)

func RunStorageBackendRaceTests(t *testing.T, _ StorageBackendFactory) {
	t.Run("ConcurrentAddHost", func(t *testing.T) {
		t.Skip("Race tests are disabled when the 'race' build tag is not set.")
	})

	t.Run("ConcurrentGetHost", func(t *testing.T) {
		t.Skip("Race tests are disabled when the 'race' build tag is not set.")
	})

	t.Run("ConcurrentRemoveHost", func(t *testing.T) {
		t.Skip("Race tests are disabled when the 'race' build tag is not set.")
	})

	t.Run("ConcurrentGetHosts", func(t *testing.T) {
		t.Skip("Race tests are disabled when the 'race' build tag is not set.")
	})
}
