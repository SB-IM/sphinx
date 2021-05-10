package uploader

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rs/zerolog"
)

type Service struct {
	client mqtt.Client
	logger *zerolog.Logger
}

func (s *Service) Serve() error {
	return nil
}
