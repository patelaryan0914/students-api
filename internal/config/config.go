package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string
}

// struct tag
type Config struct {
	Env         string `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
	StoragePath string `yaml:"storage_path"  env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

func MustLoad() *Config{
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		flags:= flag.String("config","","Path to config file")
		flag.Parse()
		configPath =*flags

		if configPath ==""{
			log.Fatal("Config Path is not set")
		}
	}

	if _,err :=os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("Config File does not exist: %s",configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath,&cfg)
	if err != nil{
		log.Fatalf("Cannot read config file %s",err.Error())
	}
	return &cfg
}