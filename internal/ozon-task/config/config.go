package config

import (
	"errors"
	"flag"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	storeType  string
	configPath string

	ErrUnknownStoreType = errors.New("unknown store type")
)

type StoreType uint8

const (
	RedisType StoreType = iota
	PostgresType
	CacheType
)

func init() {
	flag.StringVar(&storeType, "db", "postgres", "choose the type of database")
	flag.StringVar(&configPath, "config", "config.yml", "select the configs path")
}

type Config struct {
	HTTP    *HTTP    `yaml:"http"`
	Store   *Store   `yaml:"store"`
	Service *Service `yaml:"service"`
	Grpc    *Grpc    `yaml:"grpc"`
}

type HTTP struct {
	Port string `yaml:"port"`
	Host string `yaml:"host"`
}

type (
	Store struct {
		Postgres *Postgres `yaml:"postgres"`
		Redis    *Redis    `yaml:"redis"`
	}

	Postgres struct {
		DstUrl string `yaml:"dst_url"`
	}

	Redis struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
	}

	Service struct {
		LogLevel  string `yaml:"log_level"`
		PublicUrl string `yaml:"public_url"`
		StoreType StoreType
	}

	Grpc struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	}
)

func NewConfig() (*Config, error) {
	flag.Parse()
	config := &Config{}
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	d := yaml.NewDecoder(f)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	switch storeType {
	case "postgres":
		config.Service.StoreType = PostgresType
	case "redis":
		config.Service.StoreType = RedisType
	case "cache":
		config.Service.StoreType = CacheType
	default:
		return nil, ErrUnknownStoreType
	}

	if strings.HasSuffix(config.Service.PublicUrl, "/") {
		config.Service.PublicUrl = config.Service.PublicUrl[1:]
	}

	return config, nil
}
