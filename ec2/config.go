package ec2

import (
	"fmt"
	"os"
	"strings"
)

const (
	configPrefix = "AUTOMAGICAL_INSTANCE"
)

var (
	configKeys = [...]string{"table"}
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
	return os.Getenv(fmt.Sprintf("%s_%s", configPrefix, strings.ToUpper(name)))
}
