package config

type ErrUnknownLogLevel struct {
	Level string
}

func (e *ErrUnknownLogLevel) Error() string {
	return "Unknown log level " + e.Level
}

type ErrIncompleteSSLConfig struct{}

func (e ErrIncompleteSSLConfig) Error() string {
	return "SSL is enabled but certificate and/or private key are missing"
}

type ErrMissingMAC struct{}

func (e ErrMissingMAC) Error() string {
	return "At least one host has a missing mac address."
}
