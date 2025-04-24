package testsuite

import (
	"testing"

	"github.com/heathcliff26/go-wol/pkg/server/storage/types"
	"github.com/stretchr/testify/require"
)

type StorageBackendFactory func(t *testing.T, name string) types.StorageBackend

var testHosts = []types.Host{
	{
		MAC:  "AA:BB:CC:DD:EE:FF",
		Name: "TestHost1",
	},
	{
		MAC:  "11:22:33:44:55:66",
		Name: "TestHost2",
	},
	{
		MAC:  "77:88:99:AA:BB:CC",
		Name: "TestHost3",
	},
	{
		MAC:  "FF:88:99:AA:BB:CC",
		Name: "TestHost4",
	},
	{
		MAC:  "FE:11:99:AA:BB:CC",
		Name: "TestHost5",
	},
}

func addHosts(t *testing.T, backend types.StorageBackend) {
	t.Helper()

	for _, host := range testHosts {
		err := backend.AddHost(host.MAC, host.Name)
		require.NoError(t, err, "AddHost failed for %s", host.Name)
	}
}
