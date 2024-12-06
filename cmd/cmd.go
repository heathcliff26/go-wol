package main

import (
	"github.com/heathcliff26/go-wol/pkg/server"
	"github.com/heathcliff26/go-wol/pkg/version"
	"github.com/heathcliff26/go-wol/pkg/wol"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cobra.AddTemplateFunc(
		"ProgramName", func() string {
			return version.Name
		},
	)

	rootCmd := &cobra.Command{
		Use:   version.Name,
		Short: version.Name + " power on other devices on the network via Wake-on-Lan",
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	wolCMD := wol.NewCommand()

	rootCmd.AddCommand(
		wolCMD,
		server.NewCommand(),
		version.NewCommand(),
	)

	return rootCmd
}
