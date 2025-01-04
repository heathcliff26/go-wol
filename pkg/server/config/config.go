package config

import (
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_CONFIG_PATH           = "/etc/go-wol/config.yaml"
	DEFAULT_CONFIG_PATH_CONTAINER = "/config/config.yaml"

	DEFAULT_LOG_LEVEL   = "info"
	DEFAULT_SERVER_PORT = 8080
)

var logLevel *slog.LevelVar

// Initialize the logger
func init() {
	logLevel = &slog.LevelVar{}
	opts := slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
	slog.SetDefault(logger)
}

type Config struct {
	LogLevel string    `yaml:"logLevel,omitempty"`
	Port     int       `yaml:"port,omitempty"`
	SSL      SSLConfig `yaml:"ssl,omitempty"`
	Hosts    []Host    `yaml:"hosts,omitempty"`
}

type SSLConfig struct {
	Enabled bool   `yaml:"enabled,omitempty"`
	Cert    string `yaml:"cert,omitempty"`
	Key     string `yaml:"key,omitempty"`
}

type Host struct {
	Name string `yaml:"name"`
	MAC  string `yaml:"mac"`
}

// Returns a Config with default values set
func DefaultConfig() Config {
	return Config{
		LogLevel: DEFAULT_LOG_LEVEL,
		Port:     DEFAULT_SERVER_PORT,
	}
}

// Return the path to the cobfig file
func getPath(path string) string {
	if path != "" {
		return path
	}
	if _, ok := os.LookupEnv("container"); ok {
		return DEFAULT_CONFIG_PATH_CONTAINER
	} else {
		return DEFAULT_CONFIG_PATH
	}
}

// Loads config from file, returns error if config is invalid
// Arguments:
//
//		path: Path to config file, if empty will use either DEFAULT_CONFIG_PATH or DEFAULT_CONFIG_PATH_CONTAINER
//		env: Determines if enviroment variables in the file will be expanded before decoding
//	 logLevelOverride: Override the log level given by the config
func LoadConfig(path string, env bool, logLevelOverride string) (Config, error) {
	c, err := loadConfigFile(path, env)
	if err != nil {
		return Config{}, err
	}

	if logLevelOverride == "" {
		err = setLogLevel(c.LogLevel)
	} else {
		err = setLogLevel(logLevelOverride)
	}
	if err != nil {
		return Config{}, err
	}

	if c.SSL.Enabled && (c.SSL.Cert == "" || c.SSL.Key == "") {
		return Config{}, ErrIncompleteSSLConfig{}
	}

	for _, host := range c.Hosts {
		if host.MAC == "" {
			return Config{}, ErrMissingMAC{}
		}
	}

	return c, nil
}

func loadConfigFile(path string, env bool) (Config, error) {
	c := DefaultConfig()

	p := getPath(path)

	f, err := os.ReadFile(p)
	if path == "" && os.IsNotExist(err) {
		slog.Info("No config file specified and default file does not exist, falling back to default values.", slog.String("default-path", p))
		return c, nil
	} else if err != nil {
		return Config{}, err
	}

	if env {
		f = []byte(os.ExpandEnv(string(f)))
	}

	err = yaml.Unmarshal(f, &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}

// Parse a given string and set the resulting log level
func setLogLevel(level string) error {
	switch strings.ToLower(level) {
	case "debug":
		logLevel.Set(slog.LevelDebug)
	case "info":
		logLevel.Set(slog.LevelInfo)
	case "warn":
		logLevel.Set(slog.LevelWarn)
	case "error":
		logLevel.Set(slog.LevelError)
	default:
		return &ErrUnknownLogLevel{level}
	}
	return nil
}
