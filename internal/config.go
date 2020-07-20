package internal

import (
	util "github.com/Floor-Gang/utilpkg"
	"log"
	"strings"
	"time"
)

type Config struct {
	Prefix              string        `yaml:"prefix"`      // command prefix
	NotificationChannel string        `yaml:"channel"`     // channel to report to
	Guild               string        `yaml:"guild"`       // guild to listen to
	Roles               []string      `yaml:"roles"`       // minimum required roles
	Auth                string        `yaml:"auth_server"` // auth server (github.com/authserver)
	Intervals           int           `yaml:"intervals"`   // minute-intervals to check the server
	Delay               time.Duration `yaml:"kick_delay"`  // How many seconds to wait until to automatically kick them
	UpVote              string        `yaml:"up_vote"`     // kick & inform
	DownVote            string        `yaml:"down_vote"`   // discard
}

// This will get the current configuration file. If it doesn't exist then a
// new one will be made.
func GetConfig(location string) (config Config) {
	config = Config{
		Prefix:              ".subwatch",
		NotificationChannel: "",
		Roles:               []string{"1", "2", "3", "4"},
		Auth:                "",
		Guild:               "",
		UpVote:              "",
		DownVote:            "",
		Intervals:           5,
		Delay:               10,
	}
	err := util.GetConfig(location, &config)

	if err != nil && strings.Contains(err.Error(), "default configuration") {
		log.Fatalln("A default configuration has been made.")
	}

	return config
}
