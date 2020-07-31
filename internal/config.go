package internal

import (
	util "github.com/Floor-Gang/utilpkg/config"
	"log"
	"time"
)

type Config struct {
	Token               string        `yaml:"token"`
	Prefix              string        `yaml:"prefix"`
	NotificationChannel string        `yaml:"channel"`
	Guild               string        `yaml:"guild"`
	Roles               []string      `yaml:"roles"`
	Auth                string        `yaml:"auth_server"`
	Intervals           int           `yaml:"intervals"`
	Delay               time.Duration `yaml:"kick_delay"`
	UpVote              string        `yaml:"up_vote"`
	DownVote            string        `yaml:"down_vote"`
}

// This will get the current configuration file. If it doesn't exist then a
// new one will be made.
func GetConfig(location string) (config Config) {
	config = Config{
		Prefix:    ".subwatch",
		Roles:     []string{"1", "2", "3", "4"},
		Intervals: 5,
		Delay:     10,
	}
	err := util.GetConfig(location, &config)

	if err != nil {
		log.Fatalln(err)
	}

	return config
}
