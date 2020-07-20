package internal

import (
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (bot *Bot) onMessage(_ *dg.Session, message *dg.MessageCreate) {
	if len(message.GuildID) == 0 || !strings.HasPrefix(message.Content, bot.config.Prefix) {
		return
	}

	args := strings.Split(message.Content, " ")

	if len(args) < 3 {
		return
	}

	auth, err := bot.auth.Auth(message.Author.ID)

	if err != nil {
		log.Printf(
			"Failed to authentication \"%s\", because\n%s",
			message.Author.ID, err.Error(),
		)
		return
	}

	switch args[1] {
	case "add":
		if auth.IsAdmin {
			response := bot.addRole(args[2], message.GuildID)
			_, _ = util.Reply(bot.client, message.Message, response)
		} else {
			_, _ = util.Reply(
				bot.client,
				message.Message,
				"You don't have permissions to run this command.",
			)
		}
		break
	case "delete":
		if auth.IsAdmin {
			response := bot.removeRole(args[2])
			_, _ = util.Reply(bot.client, message.Message, response)
		} else {
			_, _ = util.Reply(
				bot.client,
				message.Message,
				"You don't have permissions to run this command.",
			)
		}
		break
	}
}

func (bot *Bot) onReaction(_ *dg.Session, reaction *dg.MessageReactionAdd) {
	if report, isOK := bot.reports[reaction.MessageID]; isOK {
		// ignore bot
		if reaction.UserID == bot.client.State.User.ID {
			return
		}

		// ignore emojis that aren't upvote or downvote
		if reaction.Emoji.ID != bot.config.DownVote {
			return
		}

		report.Cancel <- true
	}
}
