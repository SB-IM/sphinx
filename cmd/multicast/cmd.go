package multicast

import (
	"github.com/SB-IM/logging"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/SB-IM/sphinx/internal/multicast"
)

const configFlagName = "config"

// Command returns a livestream command.
func Command() *cli.Command {
	var (
		logger                 zerolog.Logger
		multicastConfigOptions multicast.ConfigOptions
	)

	flags := func() (flags []cli.Flag) {
		for _, v := range [][]cli.Flag{
			loadConfigFlag(),
			multicastFlags(&multicastConfigOptions),
		} {
			flags = append(flags, v...)
		}
		return
	}()

	return &cli.Command{
		Name:  "multicast",
		Usage: "multicast consumes udp stream and casts it to multiple targets",
		Flags: flags,
		Before: func(c *cli.Context) error {
			if err := altsrc.InitInputSourceWithContext(
				flags,
				altsrc.NewTomlSourceFromFlagFunc(configFlagName),
			)(c); err != nil {
				return err
			}

			// Set up logger.
			debug := c.Bool("debug")
			logging.SetDebugMod(debug)
			logger = log.With().Str("service", "sphinx").Str("command", "multicast").Logger()

			return nil
		},
		Action: func(c *cli.Context) error {
			return multicast.Run(multicastConfigOptions, &logger)
		},
		After: func(c *cli.Context) error {
			logger.Info().Msg("exits")
			return nil
		},
	}
}

// loadConfigFlag sets a config file path for app command.
// Note: you can't set any other flags' `Required` value to `true`,
// As it conflicts with this flag. You can set only either this flag or specifically the other flags but not both.
func loadConfigFlag() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        configFlagName,
			Aliases:     []string{"c"},
			Usage:       "Config file path",
			Value:       "config/config.toml",
			DefaultText: "config/config.toml",
		},
	}
}

func multicastFlags(multicastConfigOptions *multicast.ConfigOptions) []cli.Flag {
	return []cli.Flag{
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "multicast.host",
			Usage:       "Multicast server host",
			Value:       "0.0.0.0",
			DefaultText: "0.0.0.0",
			Destination: &multicastConfigOptions.Host,
		}),
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:        "multicast.port",
			Usage:       "Multicast server port",
			Value:       5004,
			DefaultText: "5004",
			Destination: &multicastConfigOptions.Port,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "multicast.sinks",
			Usage:       "Multicast targets of udp sink addresses, it's a comma separated list, each address is composed of \"host:port\" format",
			Value:       "",
			DefaultText: "",
			Destination: &multicastConfigOptions.UDPSinkAddresses,
		}),
	}
}
