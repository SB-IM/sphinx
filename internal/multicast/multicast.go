package multicast

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type ConfigOptions struct {
	Host             string
	Port             int
	UDPSinkAddresses string
}

func Run(config ConfigOptions, logger *zerolog.Logger) error {
	addr := config.Host + ":" + strconv.Itoa(config.Port)
	return multicast(addr, config.UDPSinkAddresses, logger)
}

// multicast starts a UDP listener and dispatches data to multicast addresses.
func multicast(address, sinks string, logger *zerolog.Logger) error {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("could not resolve udp address: %w", err)
	}

	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return fmt.Errorf("listen UDP: %w", err)
	}
	logger.Info().Str("address", udpAddr.String()).Msg("UDP server started")
	defer listener.Close()

	destinations, err := parseUDPSinkAddresses(sinks)
	if err != nil {
		return fmt.Errorf("could not parse udp sink addresses %s: %w", sinks, err)
	}

	inboundRTPPacket := make([]byte, 1600) // UDP MTU
	for {
		n, _, err := listener.ReadFrom(inboundRTPPacket)
		if err != nil {
			return fmt.Errorf("error during read: %w", err)
		}

		// Broadcast udp stream, acts like gstreamer multiudpsink.
		for _, dest := range destinations {
			if _, err := listener.WriteTo(inboundRTPPacket[:n], dest); err != nil {
				logger.Err(err).Str("dest", dest.String()).Msg("could not write udp packet to destination")
				// Don't return even if error happens.
			}
		}
	}
}

// parseUDPSinkAddresses parses a string list of udp addresses list, it takes the form of "192.0.0.1:2000,192.0.0.2:2001".
func parseUDPSinkAddresses(udpSinkAddresses string) (udpAddrs []*net.UDPAddr, err error) {
	if udpSinkAddresses == "" {
		return nil, errors.New("empty udpSinkAddresses")
	}
	addresses := strings.Split(udpSinkAddresses, ",")

	for _, addr := range addresses {
		udpAddr, err := net.ResolveUDPAddr("udp", addr)
		if err != nil {
			return nil, fmt.Errorf("could not resolve udp address: %w", err)
		}
		udpAddrs = append(udpAddrs, udpAddr)
	}
	return
}
