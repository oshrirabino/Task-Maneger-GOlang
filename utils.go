package main

import (
	"encoding/json"
	"os"
)

func FlagIndex(args []string) int {
	for i, a := range args {
		if len(a) > 2 && a[:2] == "--" {
			return i
		}
	}
	return -1
}

func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
