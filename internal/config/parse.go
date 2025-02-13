package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseConfig(path string) (cfg *Config, err error) {
	fp, err := os.Open(path)
	if err != nil {
		fmt.Println("Opening configuration file:", err)
		return
	}

	defer func() {
		err = fp.Close()
		if err != nil {
			fmt.Println("Error when try fp.Close[ParseConfig], err:", err)
		}
	}()

	cfg = new(Config)
	cfg.SetDefaults()
	decoder := yaml.NewDecoder(fp)
	if err = decoder.Decode(cfg); err != nil {
		fmt.Println("Invalid configuration YAML:", err)
		return
	}

	return
}
