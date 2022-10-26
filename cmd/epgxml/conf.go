package main

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Cfg struct {
	Db struct {
		User     string `toml:"user"`
		Password string `toml:"password"`
		Dbpath   string `toml:"dbpath"`
	} `toml:"db"`

	Xml struct {
		FileName string `toml:"filename"`
	} `toml:"xml"`
}

func LoadConfig(confPath string) (*Cfg, error) {
	conf := new(Cfg)
	file, err := os.Open(confPath)
	defer func() {
		if err := file.Close(); err != nil {
			LOG.Panic(err)
		}
	}()
	if err != nil {
		return nil, err
	}

	_, err = toml.DecodeFile(confPath, &conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}
