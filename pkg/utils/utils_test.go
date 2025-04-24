package utils

import (
	"strings"
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

func TestValidateHostname(t *testing.T) {
	tMatrix := []struct {
		name     string
		hostname string
		expected bool
	}{
		{"ValidHostnameSimple", "example.com", true},
		{"ValidHostnameWithSubdomain", "sub.example.com", true},
		{"ValidHostnameSingleLabelMaxLength", strings.Repeat("a", 63) + ".com", true},
		{"ValidHostnameWithNumbers", "123example.com", true},
		{"ValidHostnameWithHyphen", "example-site.com", true},
		{"ValidHostnameSingleLabel", "localhost", true},
		{"ValidHostnameWithNumericTLD", "example.123", true},
		{"ValidHostnameWithMixedCase", "ExAmPlE.cOm", true},
		{"ValidHostnameWithMultipleSubdomains", "sub1.sub2.example.com", true},
		{"ValidHostnameWithNumericSubdomain", "123.example.com", true},
		{"ValidHostnameWithMaxLength", strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 61), true},
		{"ValidSingleCharacterHostname", "a", true},
		{"InvalidHostnameSingleLabelTooLong", strings.Repeat("a", 64) + ".com", false},
		{"InvalidHostnameEmpty", "", false},
		{"InvalidHostnameWithSpecialChars", "example@site.com", false},
		{"InvalidHostnameWithUnderscore", "example_site.com", false},
		{"InvalidHostnameStartingWithHyphen", "-example.com", false},
		{"InvalidHostnameEndingWithHyphen", "example-.com", false},
		{"InvalidHostnameWithConsecutiveDots", "example..com", false},
		{"InvalidHostnameWithSpaces", "example site.com", false},
		{"InvalidHostnameWithUnicode", "exämple.com", false},
		{"InvalidHostnameWithEmoji", "example😊.com", false},
		{"InvalidHostnameTrailingDot", "example.com.", false},
		{"InvalidHostnameLeadingDot", ".example.com", false},
		{"InvalidHostnameTooLong", strings.Repeat("a", 254), false},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.name, func(t *testing.T) {
			assert := assert.New(t)

			result := ValidateHostname(tCase.hostname)

			if tCase.expected {
				assert.True(result, "ValidateHostname(%q) should be valid", tCase.hostname)
			} else {
				assert.False(result, "ValidateHostname(%q) should be invalid", tCase.hostname)
			}
		})
	}
}
