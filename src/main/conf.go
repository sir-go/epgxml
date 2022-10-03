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

func LoadConfig(confpath string) (*Cfg, error) {
	conf := new(Cfg)
	file, err := os.Open(confpath)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	_, err = toml.DecodeFile(confpath, &conf)
	if err != nil {
		return nil, err
	}

	return conf, err
}
