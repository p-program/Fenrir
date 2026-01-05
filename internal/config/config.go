package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Database DatabaseConfig `json:",optional"`
}

type DatabaseConfig struct {
	Type string `json:",default=sqlite"`
	DSN  string `json:",default=restaurant.db"`
}
