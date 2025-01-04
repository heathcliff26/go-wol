package config

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidConfigs(t *testing.T) {
	tMatrix := []struct {
		Name, Path string
		Result     Config
	}{
		{
			Name:   "EmptyConfig",
			Path:   "",
			Result: DefaultConfig(),
		},
		{
			Name: "ValidConfig1",
			Path: "testdata/valid-config-1.yaml",
			Result: Config{
				LogLevel: "debug",
				Port:     1234,
			},
		},
		{
			Name: "ValidConfig2",
			Path: "testdata/valid-config-2.yaml",
			Result: Config{
				LogLevel: "error",
				Port:     5678,
				SSL: SSLConfig{
					Enabled: true,
					Key:     "test.key",
					Cert:    "test.crt",
				},
			},
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			c, err := LoadConfig(tCase.Path, false, "")

			assert := assert.New(t)

			if !assert.NoError(err) {
				t.FailNow()
			}
			assert.Equal(tCase.Result, c)
		})
	}
}

func TestInvalidConfig(t *testing.T) {
	tMatrix := []struct {
		Name, Path, Error string
	}{
		{
			Name:  "NotYaml",
			Path:  "testdata/not-a-config.txt",
			Error: "*yaml.TypeError",
		},
		{
			Name:  "FileDoesNotExist",
			Path:  "not-a-file",
			Error: "*fs.PathError",
		},
		{
			Name:  "InvalidLogLevel",
			Path:  "testdata/invalid-config-loglevel.yaml",
			Error: "*config.ErrUnknownLogLevel",
		},
		{
			Name:  "ServerIncompleteSSLConfig1",
			Path:  "testdata/invalid-config-ssl-1.yaml",
			Error: "config.ErrIncompleteSSLConfig",
		},
		{
			Name:  "ServerIncompleteSSLConfig2",
			Path:  "testdata/invalid-config-ssl-2.yaml",
			Error: "config.ErrIncompleteSSLConfig",
		},
		{
			Name:  "HostWithEmptyMAC",
			Path:  "testdata/invalid-config-mac.yaml",
			Error: "config.ErrMissingMAC",
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			_, err := LoadConfig(tCase.Path, false, "")

			if !assert.Error(t, err) {
				t.Fatal("Did not receive an error")
			}
			if !assert.Equal(t, tCase.Error, reflect.TypeOf(err).String()) {
				t.Fatalf("Received invalid error: %v", err)
			}
		})
	}
}

func TestEnvSubstitution(t *testing.T) {
	tMatrix := []struct {
		Name   string
		Env    bool
		Config Config
	}{
		{
			Name: "Enabled",
			Env:  true,
			Config: Config{
				LogLevel: "debug",
				Port:     1234,
			},
		},
		{
			Name: "Disabled",
			Env:  false,
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			assert := assert.New(t)

			t.Setenv("GOWOL_LOG_LEVEL", "debug")
			t.Setenv("GOWOL_PORT", "1234")

			c, err := LoadConfig("testdata/env-config.yaml", tCase.Env, "")
			if tCase.Env {
				assert.NoError(err)
				assert.Equal(tCase.Config, c)
			} else {
				assert.Error(err)
			}
		})
	}
}

func TestLogLevelOverride(t *testing.T) {
	assert := assert.New(t)

	_, err := LoadConfig("testdata/valid-config-1.yaml", false, "")
	assert.NoError(err)
	assert.Equal(logLevel.Level(), slog.LevelDebug, "Should use level from config")

	_, err = LoadConfig("testdata/valid-config-1.yaml", false, "warn")
	assert.NoError(err)
	assert.Equal(logLevel.Level(), slog.LevelWarn, "Should override the log level")
}

func TestGetPath(t *testing.T) {
	tMatrix := []struct {
		Name, Path string
		Container  bool
		Result     string
	}{
		{
			Name:   "GivenPath",
			Path:   "testpath",
			Result: "testpath",
		},
		{
			Name:   "Default",
			Result: DEFAULT_CONFIG_PATH,
		},
		{
			Name:      "Container",
			Container: true,
			Result:    DEFAULT_CONFIG_PATH_CONTAINER,
		},
	}

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			if tCase.Container {
				t.Setenv("container", "podman")
			}
			assert.Equal(t, tCase.Result, getPath(tCase.Path))
		})
	}
}

func TestSetLogLevel(t *testing.T) {
	tMatrix := []struct {
		Name  string
		Level slog.Level
		Error error
	}{
		{"debug", slog.LevelDebug, nil},
		{"info", slog.LevelInfo, nil},
		{"warn", slog.LevelWarn, nil},
		{"error", slog.LevelError, nil},
		{"DEBUG", slog.LevelDebug, nil},
		{"INFO", slog.LevelInfo, nil},
		{"WARN", slog.LevelWarn, nil},
		{"ERROR", slog.LevelError, nil},
		{"Unknown", 0, &ErrUnknownLogLevel{"Unknown"}},
	}
	t.Cleanup(func() {
		err := setLogLevel(DEFAULT_LOG_LEVEL)
		if err != nil {
			t.Fatalf("Failed to cleanup after test: %v", err)
		}
	})

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			err := setLogLevel(tCase.Name)

			assert := assert.New(t)

			if !assert.Equal(tCase.Error, err) {
				t.Fatalf("Received invalid error: %v", err)
			}
			if err == nil {
				assert.Equal(tCase.Level, logLevel.Level())
			}
		})
	}
}
