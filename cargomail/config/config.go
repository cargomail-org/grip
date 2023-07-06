package config

import (
	_ "embed"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

//go:embed default.yaml
var defaultConfig []byte

type StartFlags = struct {
	DomainName   string `yaml:"domain_name"`
	DbPath       string `yaml:"db_path"`
	StoragePath  string `yaml:"storage_path"`
	AppApiBind   string `yaml:"app_api_bind"`
	GripApiBind  string `yaml:"grip_api_bind"`
	GripCertFile string `yaml:"grip_cert_file"`
	GripKeyFile  string `yaml:"grip_key_file"`
}

func NewStartFlags() StartFlags {
	sf := StartFlags{}

	setDefaults(&sf)
	loadConfig(&sf)

	return sf
}

func setDefaults(sf *StartFlags) {
	err := yaml.Unmarshal(defaultConfig, sf)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < reflect.TypeOf(*sf).NumField(); i++ {
		field := reflect.TypeOf(*sf).Field(i)
		if value, ok := field.Tag.Lookup("yaml"); ok {
			reflect.ValueOf(sf).Elem().FieldByName(field.Name).Set(reflect.ValueOf(os.Getenv(strings.ToUpper(value))))
		}
	}
}

func loadConfig(sf *StartFlags) {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(configFile, sf)
	if err != nil {
		log.Fatal(err)
	}
}
