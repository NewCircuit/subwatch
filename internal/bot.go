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
	client.AddHandlerOnce(bot.onReady)
	client.AddHandler(bot.onMessage)

	if err := client.Open(); err != nil {
		log.Fatalln("Failed to connect to Discord, was an access token provided?\n" + err.Error())
	}
}

func (b *Bot) onReady(_ *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Ready as %s#%s\n", ready.User.Username, ready.User.Discriminator)
	b.reviewGuild()	
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <- ticker.C:
				b.reviewGuild()
			}
		}
	}()

}

func (b *Bot) reviewGuild() {
	result := "**__SubWatch__** ⚠\nThese people need to be checked up on:\n"
	members := ""
	b.reviewMembers("", &members)

	if len(members) > 0 {
		_, err := b.client.ChannelMessageSend(b.config.NotificationChannel, result + members)
		if err != nil {
			log.Printf("Failed to send a report to \"%s\" because\n%s\n", b.config.NotificationChannel, err.Error())
		}
	}
}
 
func (b *Bot) reviewMembers(memberID string, result *string) {
	members, err := b.client.GuildMembers(b.config.Guild, "", 1000)
	var last string

	if err != nil {
		log.Printf("Failed to fetch members for \"%s\" because\n%s\n", b.config.Guild, err.Error())
		return
	}

	for _, member := range members {
		if !b.checkRoles(member.Roles) {
			*result += fmt.Sprintf(" - %s#%s (<@%s>)\n", member.User.Username, member.User.Discriminator, member.User.ID)
		}
		last = member.User.ID
	}

	if len(members) == 1000 {
		b.reviewMembers(last, result)
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

func (b Bot) reply(event *discordgo.MessageCreate, context string) (*discordgo.Message, error) {
	return b.client.ChannelMessageSend(
		event.ChannelID,
		fmt.Sprintf("<@%s> %s", event.Author.ID, context),
	)
}
