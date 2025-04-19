package server

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/heathcliff26/go-wol/pkg/server/config"
	"github.com/spf13/cobra"
)

const (
	flagNameConfig   = "config"
	flagNameLogLevel = "log"
	flagNameEnv      = "env"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Serve a frontend via gui",
		Run: func(cmd *cobra.Command, args []string) {
			err := run(cmd)
			if err != nil {
				fmt.Println("Fatal: " + err.Error())
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringP(flagNameConfig, "c", "", "Config file to use")
	cmd.Flags().String(flagNameLogLevel, "", "Override the log level given in the config file")
	cmd.Flags().Bool(flagNameEnv, false, "Expand enviroment variables in the config file")

	return cmd
}

func run(cmd *cobra.Command) error {
	configPath, err := cmd.Flags().GetString(flagNameConfig)
	if err != nil {
		return err
	}
	logLevel, err := cmd.Flags().GetString(flagNameLogLevel)
	if err != nil {
		return err
	}
	env, err := cmd.Flags().GetBool(flagNameEnv)
	if err != nil {
		return err
	}

	cfg, err := config.LoadConfig(configPath, env, logLevel)
	if err != nil {
		return err
	}

	// TODO: Remove when api functions have been implemented
	if !cfg.Storage.Readonly {
		slog.Warn("Writing to storage is not yet supported, ignoring setting")
	}
	cfg.Storage.Readonly = true

	server, err := NewServer(cfg.Server, cfg.Storage)
	if err != nil {
		return err
	}

	return server.Run()
}
