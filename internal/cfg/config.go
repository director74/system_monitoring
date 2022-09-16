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
}

type Config struct {
	TrackAllowed TrackAllowedConf
}

type TrackAllowedConf struct {
	LoadAverage bool `yaml:"loadAverage"`
	CpuLoad     bool `yaml:"  cpuLoad"`
	DiskLoad    bool `yaml:"diskLoad"`
	TopTalkers  bool `yaml:"topTalkers"`
	NetStats    bool `yaml:"netStats"`
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
