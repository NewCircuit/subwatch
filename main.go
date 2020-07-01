package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/floor-gang/RoleWatcher/helpers"
	"os"
	"os/signal"
	"syscall"
)

const (
	configLocation = "config.yml"
)

type Bot struct {
	config helpers.Config
	client *discordgo.Session
}

func main() {
	config := helpers.GetConfig(configLocation)
	client, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		panic(err)
	}

	bot := Bot{
		config: config,
		client: client,
	}

	client.AddHandler(bot.onRoleUpdate)
	client.AddHandler(bot.onMessage)

	if err = client.Open(); err != nil {
		panic(err)
	}

	keepAlive()
}

func keepAlive() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}