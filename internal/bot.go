package internal

import (
	"fmt"
	auth "github.com/Floor-Gang/authclient"
	"github.com/bwmarrin/discordgo"
	"log"
)

type Bot struct {
	config  Config
	client  *discordgo.Session
	confLoc string
	auth    auth.AuthClient
	reports map[string]Report
}

type Report struct {
	MemberIDs []string
	ReportID  string
	Cancel    chan bool
}

func Start(config Config, configLocation string) {
	// Setup Authentication client
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to authentication server because \n" + err.Error())
	}

	_, err = authClient.Register(
		auth.Feature{
			Name:        "SubWatch",
			Description: "This is responsible for reporting users that aren't paying anymore.",
			Commands: []auth.SubCommand{
				{
					Name:        "add",
					Description: "Add another role",
					Example:     []string{"add", "role ID"},
				},
				{
					Name:        "delete",
					Description: "Remove a role",
					Example:     []string{"remove", "role ID"},
				},
			},
			CommandPrefix: config.Prefix,
		},
	)
	if err != nil {
		log.Fatalln("Failed to register with the authentication server\n" + err.Error())
	}

	// Setup Discord
	client, _ := discordgo.New("Bot " + config.Token)

	// Declare intents
	intents := discordgo.MakeIntent(
		discordgo.IntentsGuildMembers + discordgo.IntentsGuildMessages +
			discordgo.IntentsGuildMessageReactions,
	)
	client.Identify.Intents = intents

	bot := Bot{
		config:  config,
		client:  client,
		confLoc: configLocation,
		auth:    authClient,
		reports: map[string]Report{},
	}

	// Add event listeners
	client.AddHandlerOnce(bot.onReady)
	client.AddHandler(bot.onMessage)
	client.AddHandler(bot.onReaction)

	// connect to Discord websocket
	if err := client.Open(); err != nil {
		log.Fatalln(
			"Failed to connect to Discord, was an access token provided?\n" +
				err.Error(),
		)
	}
}

// check if they have at least one of the required roles from the config.
func (bot Bot) checkRoles(userRoles []string) bool {
	for _, role := range userRoles {
		if hasRole(role, bot.config.Roles) {
			return true
		}
	}

	return false
}

// kick all the given member IDs
func (bot *Bot) kickMembers(members []string) string {
	var result = "**SubWatch Report - Conclusion**\n"

	for _, memberID := range members {
		// first DM the member
		dmChannel, dmErr := bot.client.UserChannelCreate(memberID)
		user, _ := bot.client.User(memberID)

		if dmErr == nil {
			_, dmErr = bot.client.ChannelMessageSend(
				dmChannel.ID,
				"Renew your membership",
			)
		}
		// then kick them
		kickErr := bot.client.GuildMemberDeleteWithReason(
			bot.config.Guild,
			memberID,
			"Renew your membership",
		)

		// if we failed to contact then let them know
		if dmErr != nil {
			result += fmt.Sprintf(" - %s#%s (failed to notify & ", user.Username, user.Discriminator)
		} else {
			result += fmt.Sprintf(" - %s#%s (notified & ", user.Username, user.Discriminator)
		}

		if kickErr == nil && dmErr != nil {
			result += "kicked anyways)"
		} else if kickErr == nil && dmErr == nil {
			result += "kicked)"
		} else {
			result += "failed to kick do I have the required perms?)"
		}
		result += "\n"

	}
	return result
}
