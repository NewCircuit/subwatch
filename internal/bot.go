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
}

func Start(config Config, configLocation string) {
	// Setup Authentication client
	authClient, err := auth.GetClient(config.Auth)

	if err != nil {
		log.Fatalln("Failed to connect to authentication server because \n" + err.Error())
	}

	register, err := authClient.Register(
		auth.Feature{
			Name:        "SubWatch",
			Description: "This is responsible for reporting users that aren't paying anymore.",
			Commands: []auth.SubCommand{
				{
					Name:        "add",
					Description: "Add another role",
					Example:     []string{"add", "<role ID>"},
				},
				{
					Name:        "delete",
					Description: "Remove a role",
					Example:     []string{"remove", "<role ID>"},
				},
			},
			CommandPrefix: config.Prefix,
		},
	)

	if err != nil {
		log.Fatalln("Failed to register with the authenticaiton server\n" + err.Error())
	}

	// Setup Discord
	client, _ := discordgo.New(register.Token)

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
func (b Bot) checkRoles(userRoles []string) bool {
	for _, role := range userRoles {
		if hasRole(role, b.config.Roles) {
			return true
		}
	}

	return false
}

// kick all the given member IDs
func (b *Bot) kickMembers(members []string) (failures []string) {
	for _, memberID := range members {
		// first DM the member
		dmChannel, err := b.client.UserChannelCreate(memberID)

		if err == nil {
			_, err = b.client.ChannelMessageSend(
				dmChannel.ID,
				"Renew your membership",
			)
		}

		// if we failed to contact them then don't kick them.
		if err != nil {
			failures = append(
				failures,
				fmt.Sprintf("Couldn't contact <@%s>", memberID),
			)
			continue
		}

		// then kick them if they were dm'd
		err = b.client.GuildMemberDeleteWithReason(
			b.config.Guild,
			memberID,
			"Renew your membership",
		)

		if err != nil {
			failures = append(
				failures,
				fmt.Sprintf(
					"Couldn't kick <@%s>, do I have the perms?",
					memberID,
				),
			)
		}
	}
	return failures
}
