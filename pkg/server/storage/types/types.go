package types

// Host on the network.
type Host struct {
	MAC  string `json:"mac"`
	Name string `json:"name"`
}

// StorageBackend is a concurrency safe interface for different methods of storing the configured hosts.
type StorageBackend interface {
	// Add a new host, overwrite existing host name if it already exists.
	// Ensures that the MAC address is unique and uppercase.
	AddHost(mac, host string) error
	// Remove a host, ignore if the host does not exist
	RemoveHost(mac string) error
	// Return the host name for a given MAC address, return empty if not found
	GetHost(mac string) (string, error)
	// Return all hosts
	GetHosts() ([]Host, error)
	// Check if the storage backend is readonly
	Readonly() (bool, error)
}

// Struct for reading hosts from a yaml file.
type HostsFile struct {
	Hosts []Host `json:"hosts"`
}
