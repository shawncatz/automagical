package ec2

import (
	"fmt"
	"os"
	"strings"
)

const (
	configPrefix = "AUTOMAGICAL_EC2"
)

var (
	configKeys     = [...]string{"table"}
	configDefaults = map[string]string{
		"AUTOMAGICAL_EC2_TABLE": "automagical_ec2",
	}
)

type Config map[string]string

func NewConfig() Config {
	c := Config{}
	for _, k := range configKeys {
		c[k] = env(k)
	}
	return c
}

func env(name string) string {
	key := fmt.Sprintf("%s_%s", configPrefix, strings.ToUpper(name))
	if s := os.Getenv(key); s != "" {
		return s
	}

	return configDefaults[key]
}
