package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMACAddress(t *testing.T) {
	tMatrix := []struct {
		name     string
		macAddr  string
		expected bool
	}{
		{"ValidMacAddress", "01:23:45:67:89:AB", true},
		{"InvalidMacAddressTooShort", "01:23:45:67:89", false},
		{"InvalidMacAddressInvalidCharacters", "01:23:45:67:89:ZZ", false},
		{"ValidMacAddressLowercase", "01:23:45:67:89:ab", true},
		{"InvalidMacAddressEmptyString", "", false},
		{"InvalidMacAddressMissingColons", "0123456789AB", false},
		{"InvalidMacAddressTooLong", "01:23:45:67:89:AB:CD", false},
		{"ValidMacAddressWithDash", "01-23-45-67-89-AB", true},
		{"InvalidMacAddressLeadingSpaces", " 01:23:45:67:89:AB", false},
		{"InvalidMacAddressTrailingSpaces", "01:23:45:67:89:AB ", false},
		{"InvalidMacAddressLeadingAndTrailingSpaces", " 01:23:45:67:89:AB ", false},
		{"InvalidMacAddressMixedSeparators", "01-23:45-67:89:ab", false},
		{"InvalidMacAddressWithUnicode", "01:23:45:67:89:ÄB", false},
		{"InvalidMacAddressWithEmoji", "01:23:45:67:89:😊", false},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.name, func(t *testing.T) {
			result := ValidateMACAddress(tCase.macAddr)
			assert.Equalf(t, tCase.expected, result, "ValidateMACAddress(%q)", tCase.macAddr)
		})
	}
}
