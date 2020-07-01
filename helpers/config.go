package helpers

import (
	"fmt"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Token      string
	Prefix     string
	NotificationChannel string
	Roles []string
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

// This will create a new configuration file
func genConfig(location string) Config {
	newConfig := Config{
		Token: "",
		Prefix: ".role-watcher",
		NotificationChannel: "",
		Roles: []string{"718524998234669108", "718464650194452511", "719805582969798677", "718454001279959071", "721112261283807273", "719544792828215398", "718607495437746207", "722864056284872784", "718453523057999952", "718816943452323880"},
	}

	serialized, err := yaml.Marshal(newConfig)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(location, serialized, 0660)

	if err != nil {
		panic(err)
	}

	fmt.Println("Config has been generated, please take a look at " + location)
	return GetConfig(location)
}