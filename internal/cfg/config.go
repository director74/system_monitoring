package cfg

import (
	"fmt"
	"os"

	yml "gopkg.in/yaml.v2"
)

//go:generate mockgen -source=./internal/cfg/config.go --destination=./test/mocks/cfg/config.go
type Configurable interface {
	Parse(path string) error
	GetAllowedForTracking() TrackAllowedConf
	GetGRPCServerConf() GRPCServerConf
	GetClearPeriodConf() ClearPeriodConf
}

type Config struct {
	TrackAllowed TrackAllowedConf `yaml:"TrackAllowed"`
	GRPCServer   GRPCServerConf   `yaml:"GRPCServer"`
	ClearPeriod  ClearPeriodConf  `yaml:"ClearPeriod"`
}

type TrackAllowedConf struct {
	LoadAverage bool `yaml:"LoadAverage"`
	CpuLoad     bool `yaml:"CpuLoad"`
	DiskLoad    bool `yaml:"DiskLoad"`
	TopTalkers  bool `yaml:"TopTalkers"`
	NetStats    bool `yaml:"NetStats"`
}

type GRPCServerConf struct {
	Host string `yaml:"Host"`
	Port string `yaml:"Port"`
}

type ClearPeriodConf struct {
	Minutes int `yaml:"Minutes"`
}

func NewConfig() Configurable {
	return &Config{}
}

func (c *Config) Parse(path string) error {
	configYml, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading %v error: %w", path, err)
	}

	err = yml.Unmarshal(configYml, c)
	if err != nil {
		return fmt.Errorf("can't parse %v: %w", path, err)
	}

	return nil
}

func (c *Config) GetAllowedForTracking() TrackAllowedConf {
	return c.TrackAllowed
}

func (c *Config) GetGRPCServerConf() GRPCServerConf {
	return c.GRPCServer
}

func (c *Config) GetClearPeriodConf() ClearPeriodConf {
	return c.ClearPeriod
}
