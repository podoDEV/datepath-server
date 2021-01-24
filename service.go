package datepath_server

import (
	"github.com/PODO/datepath-server/pkg/direction"
	"github.com/PODO/datepath-server/pkg/local"
)

type ServiceClients struct {
	Local     *local.Client
	Direction *direction.Client
}

func NewServiceClients(conf *Config) *ServiceClients {
	return &ServiceClients{
		Local:     local.NewClient(&conf.Local),
		Direction: direction.NewClient(&conf.Direction),
	}
}
