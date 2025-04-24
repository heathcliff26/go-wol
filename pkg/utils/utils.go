package utils

import (
	"net"
	"regexp"
	"strings"
)

var hostnameValidCharsRegexp = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$`)

// ValidateMACAddress checks if the given string is a valid MAC address.
func ValidateMACAddress(macAddrStr string) bool {
	_, err := net.ParseMAC(macAddrStr)
	return err == nil
}

// Validate that the given hostname is a valid domain name.
func ValidateHostname(hostname string) bool {
	// Hostname must be between 1 and 253 characters
	if len(hostname) < 1 || len(hostname) > 253 {
		return false
	}

	labels := strings.Split(hostname, ".")
	for _, label := range labels {
		// Each label must be between 1 and 63 characters
		if len(label) < 1 || len(label) > 63 {
			return false
		}

		if !hostnameValidCharsRegexp.MatchString(label) {
			return false
		}
	}

	return true
}
