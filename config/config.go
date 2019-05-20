package config

import (
	base "github.com/alive2212/go-illuminate/config"
	"os"
	"time"

)

var _config *Config

//Config ...
type Config struct {
}

//ReadConfig reads config file using given path
func ReadConfig(filePath string) (*Config, error) {
	if _config != nil {
		return _config, nil
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read config file")
	}
	config := &Config{}
	err = yaml.Unmarshal(content, config)
	if err != nil {
		return nil, errors.Wrap(err, "could not unmarshall config file")
	}
	_config = config

	return config, nil
}
