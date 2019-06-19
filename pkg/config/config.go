package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nkuba/btc-tools/pkg/chain"
)

// Config is the top level config structure.
type Config struct {
	BlockCypher chain.Config
	Signer      Signer
}

// Signer holds a signer config.
type Signer struct {
	WIF string
}

// ReadConfig reads in the configuration file in .toml format.
func ReadConfig(filePath string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		return nil, fmt.Errorf("unable to decode .toml file [%s] error [%s]", filePath, err)
	}

	return config, nil
}
