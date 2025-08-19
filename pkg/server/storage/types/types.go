package types

// Host on the network.
type Host struct {
	MAC     string `json:"mac" validate:"required" example:"AA:BB:CC:DD:EE:FF"`
	Name    string `json:"name" validate:"required" example:"my-host"`
	Address string `json:"address,omitempty" validate:"optional" example:"host.example.org"`
}

// StorageBackend is a concurrency safe interface for different methods of storing the configured hosts.
type StorageBackend interface {
	// Add a new host, overwrite existing host name if it already exists.
	// Ensures that the MAC address is unique and uppercase.
	AddHost(host Host) error
	// Remove a host, ignore if the host does not exist
	RemoveHost(mac string) error
	// Return the host name for a given MAC address, return empty if not found
	GetHost(mac string) (Host, error)
	// Return all hosts
	GetHosts() ([]Host, error)
	// Check if the storage backend is readonly
	Readonly() (bool, error)
}

// Struct for reading hosts from a yaml file.
type HostsFile struct {
	Hosts []Host `json:"hosts"`
}
