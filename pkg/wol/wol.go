package wol

import (
	"encoding/binary"
	"fmt"
	"log/slog"
	"net"
)

const DEFAULT_BROADCAST_ADDRESS = "255.255.255.255"

type MACAddress [6]byte

type MagicPacket struct {
	// The header is 6 bytes of 0xFF
	header [6]byte
	// The header consists of the mac address repeated 16 times
	payload [16]MACAddress
}

// Create a new magic packet from the given mac address
func CreatePacket(macAddrStr string) (*MagicPacket, error) {
	hwAddr, err := net.ParseMAC(macAddrStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse MAC address '%s': %w", macAddrStr, err)
	}

	packet := &MagicPacket{}

	for i := range packet.header {
		packet.header[i] = 0xFF
	}

	var macAddr MACAddress

	for i := range macAddr {
		macAddr[i] = hwAddr[i]
	}

	for i := range packet.payload {
		packet.payload[i] = macAddr
	}

	return packet, nil
}

// Send the magic packet to the given broadcast address
func (p *MagicPacket) Send(bcAddr string) error {
	if bcAddr == "" {
		bcAddr = DEFAULT_BROADCAST_ADDRESS
	}

	buf, err := binary.Append(nil, binary.BigEndian, p)
	if err != nil {
		return fmt.Errorf("failed to serialize magic packet: %w", err)
	}

	conn, err := net.Dial("udp", bcAddr+":9")
	if err != nil {
		return fmt.Errorf("failed to dial UDP address '%s:9': %w", bcAddr, err)
	}
	defer conn.Close()

	bytesWritten, err := conn.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to send magic packet to '%s:9': %w", bcAddr, err)
	}

	slog.Debug("Send packet", slog.String("broadcast", bcAddr), slog.Int("bytesWritten", bytesWritten))
	return nil
}
