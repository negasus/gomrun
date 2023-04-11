package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Build struct {
	Context string `yaml:"context"`
	Path    string `yaml:"path"`
}

type Service struct {
	WorkDir     string            `yaml:"work_dir"`
	Cmd         string            `yaml:"cmd"`
	Build       *Build            `yaml:"build"`
	Args        []string          `yaml:"args"`
	Envset      []string          `yaml:"envset"`
	EnvFile     string            `yaml:"env_file"`
	Environment map[string]string `yaml:"environment"`
	Delay       int               `yaml:"delay"`
}

func (c *Service) validate(cfg *Config) error {
	if c.Build == nil && len(c.Cmd) == 0 || c.Build != nil && len(c.Cmd) > 0 {
		return fmt.Errorf("you must use only Build or Cmd section on Service, not both")
	}

	for _, n := range c.Envset {
		es, ok := cfg.Envset[n]
		if !ok || es == nil {
			return fmt.Errorf("envset %q not found", n)
		}
	}

	if c.Delay < 0 {
		return fmt.Errorf("delay must be greater or equals 0")
	}

	return nil
}

type Config struct {
	Services map[string]Service           `yaml:"services"`
	Envset   map[string]map[string]string `yaml:"envset"`
}

func Load(data []byte) (*Config, error) {
	cfg := &Config{
		Services: map[string]Service{},
		Envset:   map[string]map[string]string{},
	}

	errDecode := yaml.Unmarshal(data, cfg)
	if errDecode != nil {
		return nil, fmt.Errorf("error parse config, %w", errDecode)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("error validate config, %w", err)
	}

	return cfg, nil

}

func (c *Config) validate() error {
	if c.Services == nil {
		return fmt.Errorf("no services")
	}
	if len(c.Services) == 0 {
		return fmt.Errorf("no services")
	}

	for name, srv := range c.Services {
		if err := srv.validate(c); err != nil {
			return fmt.Errorf("service %s, %v", name, err)
		}
	}
	return nil
}
