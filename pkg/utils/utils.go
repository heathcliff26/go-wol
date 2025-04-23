package utils

import "net"

// ValidateMACAddress checks if the given string is a valid MAC address.
func ValidateMACAddress(macAddrStr string) bool {
	_, err := net.ParseMAC(macAddrStr)
	return err == nil
}
