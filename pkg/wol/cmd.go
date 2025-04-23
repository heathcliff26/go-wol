package wol

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	flagNameBroadcastAddress = "broadcast"
)

// Create new Wake-on-Lan command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wol",
		Short: "Send a magic packet to the given mac address",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			bcAddr, err := cmd.Flags().GetString(flagNameBroadcastAddress)
			if err != nil {
				exitError(cmd, err)
			}

			err = run(args[0], bcAddr)
			if err != nil {
				exitError(cmd, err)
			}
		},
	}

	cmd.Flags().StringP(flagNameBroadcastAddress, "b", DEFAULT_BROADCAST_ADDRESS, "The broadcast ip address to use")

	return cmd
}

func run(macAddress, bcAddr string) error {
	packet, err := CreatePacket(macAddress)
	if err != nil {
		return fmt.Errorf("invalid MAC address: %w", err)
	}

	err = packet.Send(bcAddr)
	if err != nil {
		return fmt.Errorf("failed to send magic packet: %w", err)
	}
	fmt.Printf("Magic packet sent to %s\n", macAddress)
	return nil
}

// Print the error information on stderr and exit with code 1
func exitError(cmd *cobra.Command, err error) {
	fmt.Fprintln(cmd.Root().ErrOrStderr(), "Fatal: "+err.Error())
	os.Exit(1)
}
