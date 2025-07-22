package valkey

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/heathcliff26/go-wol/pkg/server/storage/types"
	"github.com/valkey-io/valkey-go"
)

const hostsListKey = "hosts"

const defaultTimeout = 5 * time.Second

type ValkeyBackend struct {
	client valkey.Client
}

type ValkeyConfig struct {
	Addrs    []string `json:"addresses,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	DB       int      `json:"db,omitempty"`
	TLS      bool     `json:"tls,omitempty"`

	// Options for sentinel
	Sentinel  bool   `json:"sentinel,omitempty"`
	MasterSet string `json:"master,omitempty"`
}

func NewValkeyBackend(cfg ValkeyConfig) (*ValkeyBackend, error) {
	var client valkey.Client
	var tlsConfig *tls.Config

	if cfg.TLS {
		tlsConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	opt := valkey.ClientOption{
		InitAddress: cfg.Addrs,
		Username:    cfg.Username,
		Password:    cfg.Password,
		SelectDB:    cfg.DB,
		TLSConfig:   tlsConfig,

		DisableCache: true,
	}

	if cfg.Sentinel {
		opt.Sentinel = valkey.SentinelOption{
			MasterSet: cfg.MasterSet,
			Username:  cfg.Username,
			Password:  cfg.Password,
		}
	}

	client, err := valkey.NewClient(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to valkey server: %w", err)
	}

	return &ValkeyBackend{
		client: client,
	}, nil
}

// Add a new host, overwrite existing host name if it already exists.
// Ensures that the MAC address is unique and uppercase.
func (v *ValkeyBackend) AddHost(mac string, host string) error {
	mac = strings.ToUpper(mac)

	cmdAdd := v.client.B().Set().Key(mac).Value(host).Build()
	cmdZadd := v.client.B().Zadd().Key(hostsListKey).Nx().ScoreMember().ScoreMember(float64(time.Now().UnixNano()), mac).Build()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := v.client.Do(ctx, cmdAdd).Error()
	if err != nil {
		return fmt.Errorf("failed to set host: %w", err)
	}

	err = v.client.Do(ctx, cmdZadd).Error()
	if err != nil {
		return fmt.Errorf("failed to add host to list: %w", err)
	}

	return nil
}

// Remove a host, ignore if the host does not exist
func (v *ValkeyBackend) RemoveHost(mac string) error {
	mac = strings.ToUpper(mac)

	cmdDel := v.client.B().Del().Key(mac).Build()
	cmdZrem := v.client.B().Zrem().Key(hostsListKey).Member(mac).Build()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	err := v.client.Do(ctx, cmdZrem).Error()
	if err != nil {
		return fmt.Errorf("failed to remove host from list: %w", err)
	}

	err = v.client.Do(ctx, cmdDel).Error()
	if err != nil {
		return fmt.Errorf("failed to delete host, but already removed host from list: %w", err)
	}

	return nil
}

// Return the host name for a given MAC address, return empty if not found
func (v *ValkeyBackend) GetHost(mac string) (string, error) {
	mac = strings.ToUpper(mac)

	cmdGet := v.client.B().Get().Key(mac).Build()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	val, err := v.client.Do(ctx, cmdGet).ToString()
	if err != nil {
		return "", fmt.Errorf("failed to get host: %w", err)
	}
	return val, nil
}

// Return all hosts
func (v *ValkeyBackend) GetHosts() ([]types.Host, error) {
	cmdZrange := v.client.B().Zrange().Key(hostsListKey).Min("0").Max("-1").Build()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	macs, err := v.client.Do(ctx, cmdZrange).AsStrSlice()
	if err != nil {
		return nil, fmt.Errorf("failed to get known hosts list: %w", err)
	}

	res, err := valkey.MGet(v.client, ctx, macs)
	if err != nil {
		return nil, fmt.Errorf("failed to get host names: %w", err)
	}

	hosts := make([]types.Host, 0, len(macs))
	for _, mac := range macs {
		val, ok := res[mac]
		if !ok {
			return nil, fmt.Errorf("MAC address '%s' is in list but no hostname is found", mac)
		}

		name, err := val.ToString()
		if err != nil {
			return nil, fmt.Errorf("failed to convert response hostname value to string: %w", err)
		}

		hosts = append(hosts, types.Host{
			MAC:  mac,
			Name: name,
		})
	}
	return hosts, nil
}

// Check if the storage backend is readonly
func (v *ValkeyBackend) Readonly() (bool, error) {
	// You can configure valkey to be readonly via ACL, or connect against a replica.
	// However i do not know how to reliably check if the connection is readonly.
	// Instead of hoping that no network error occurs on startup, we just default to assume we can write.
	return false, nil
}
