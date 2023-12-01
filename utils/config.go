package utils

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

type LBStrategy int

const (
	RoundRobin LBStrategy = iota
	LeastConnection
)

func GetLBStrategy(strategy string) LBStrategy {
	switch strategy {
	case "least-connection":
		return LeastConnection
	default:
		return RoundRobin
	}
}

type config struct {
	Port            int      `yaml:"lb_port"`
	Strategy        string   `yaml:"strategy"`
	Backends        []string `yaml:"backends"`
	MaxAttemptLimit int      `yaml:"max_attempt_limit"`
}

const MAX_LB_ATTEMPT int = 3

func GetLBConfig() (*config, error) {
	var conf config
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		return nil, err
	}
	if len(conf.Backends) == 0 {
		return nil, errors.New("backend hosts expected, none provided")
	}

	if conf.Port == 0 {
		return nil, errors.New("load balancer port not found")
	}
	return &conf, nil
}
