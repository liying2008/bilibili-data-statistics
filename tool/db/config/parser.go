package config

import "encoding/json"

func ParseDBConfig(jsonData []byte) *Config {
	config := &Config{}
	json.Unmarshal(jsonData, config)
	return config
}
