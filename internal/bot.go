package internal

import (
	"fmt"
	auth "github.com/Floor-Gang/authclient"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

type Bot struct {
	config  Config
	client  *discordgo.Session
	confLoc string
	auth    auth.AuthClient
}

func Start(config Config, configLocation string) {
	// Setup Discord
	client, _ := discordgo.New("Bot " + config.Token)

	// This is required
	intents := discordgo.MakeIntent(discordgo.IntentsGuildMembers + discordgo.IntentsGuildMessages)
	client.Identify.Intents = intents

	// Setup Authentication client
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to authentication server because \n" + err.Error())
	}

	bot := Bot{
		config:  config,
		client:  client,
		confLoc: configLocation,
		auth:    authClient,
	}

	// Add event listeners
	client.AddHandler(bot.onReady)
	client.AddHandler(bot.onMemberUpdate)
	client.AddHandler(bot.onMessage)

	if err = client.Open(); err != nil {
		log.Fatalln("Failed to connect to Discord, was an access token provided?\n" + err.Error())
	}
}

func (b *Bot) onReady(_ *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Ready as %s#%s\n", ready.User.Username, ready.User.Discriminator)
}

func (b *Bot) onMemberUpdate(_ *discordgo.Session, member *discordgo.GuildMemberUpdate) {
	// Ignore other guilds
	if member.GuildID != b.config.Guild {
		return
	}
	if !b.checkRoles(member.Roles) {
		b.sendEmbed(member.Member)
	}
}

// check if they have at least one of the required roles from the config.
func (b Bot) checkRoles(userRoles []string) bool {

	for _, role := range userRoles {
		if hasRole(role, b.config.Roles) {
			return true
		}
	}

	return false
}

func (b Bot) sendEmbed(member *discordgo.Member) {
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.Username,
			IconURL: member.User.AvatarURL(""),
		},
		Color:       0xff0000,
		Description: fmt.Sprintf("<@%s> needs to be checked up on", member.User.ID),
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "Sub Watcher ⚠",
	}

	_, err := b.client.ChannelMessageSendEmbed(b.config.NotificationChannel, &embed)

	if err != nil {
		log.Printf("Failed to report %s in %s, because \n%s\n", member.User.Username, b.config.NotificationChannel, err.Error())
	} else {
		log.Printf("Reported %s in %s\n", member.User.Username, b.config.NotificationChannel)
	}
}

func (b Bot) reply(event *discordgo.MessageCreate, context string) (*discordgo.Message, error) {
	return b.client.ChannelMessageSend(
		event.ChannelID,
		fmt.Sprintf("<@%s> %s", event.Author.ID, context),
	)
}
