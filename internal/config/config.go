package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbPath   string `yaml:"db_path"`
	DumpPath string `yaml:"dump_path"`
}

func Load(confPath string) (cfg *Config, err error) {
	var b []byte
	b, err = ioutil.ReadFile(filepath.Clean(confPath))
	if err != nil {
		return
	}
	cfg = new(Config)
	err = yaml.UnmarshalStrict(b, cfg)
	return
}
