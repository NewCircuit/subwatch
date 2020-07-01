package internal

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token               string   `yaml:"token"`
	Prefix              string   `yaml:"prefix"`
	NotificationChannel string   `yaml:"channel"`
	Roles               []string `yaml:"roles"`
}

// This will get the current configuration file. If it doesn't exist then a
// new one will be made.
func GetConfig(location string) (config Config) {
	// Check if the config file exists, if it doesn't create one with a
	// template.
	if _, err := os.Stat("config.yml"); err != nil {
		genConfig(location)
		log.Fatalln("Created a default config.")
	}

	// Config file exists, so we're reading it.
	file, err := ioutil.ReadFile(location)

	if err != nil {
		log.Fatalln("Failed to read config file\n" + err.Error())
	}

	// Parse the yml file
	_ = yaml.Unmarshal(file, &config)

	return config
}

// This will create a new configuration file.
func genConfig(location string) Config {
	newConfig := Config{
		Token:               "",
		Prefix:              ".role-watcher",
		NotificationChannel: "",
		Roles:               []string{"1", "2", "3", "4"},
	}

	serialized, err := yaml.Marshal(newConfig)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(location, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return GetConfig(location)
}
