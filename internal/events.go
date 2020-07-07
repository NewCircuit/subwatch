package internal

import (
	"fmt"
	util "github.com/Floor-Gang/utilpkg"
	dg "github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

func (b *Bot) onMessage(_ *dg.Session, message *dg.MessageCreate) {
	if len(message.GuildID) == 0 || !strings.HasPrefix(message.Content, b.config.Prefix) {
		return
	}

	args := strings.Split(message.Content, " ")

	if len(args) < 3 {
		return
	}

	auth, err := b.auth.Auth(message.Author.ID)

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
			response := b.addRole(args[2], message.Author.ID, message.GuildID)
			_, _ = util.Reply(b.client, message.Message, response)
		} else {
			_, _ = util.Reply(
				b.client,
				message.Message,
				"You don't have permissions to run this command.",
			)
		}
		break
	case "delete":
		if auth.IsAdmin {
			response := b.removeRole(args[2], message.Author.ID)
			_, _ = util.Reply(b.client, message.Message, response)
		} else {
			_, _ = util.Reply(
				b.client,
				message.Message,
				"You don't have permissions to run this command.",
			)
		}
		break
	}
}

func (b *Bot) onReaction(_ *dg.Session, reaction *dg.MessageReactionAdd) {
	if report, ok := b.reports[reaction.MessageID]; ok {
		// ignore bot
		if reaction.UserID == b.client.State.User.ID {
			return
		}

		// ignore emojis that aren't upvote or downvote
		if reaction.Emoji.ID != b.config.UpVote && reaction.Emoji.ID != b.config.DownVote {
			return
		}

		toKick := reaction.Emoji.ID == b.config.UpVote

		if toKick {
			failures := b.kickMembers(report.MemberIDs)

			if len(failures) > 0 {
				failureReport := "Failures:\n"
				for _, failure := range failures {
					failureReport += fmt.Sprintf(" - %s\n", failure)
				}
				_, _ = b.client.ChannelMessageSend(
					reaction.ChannelID,
					failureReport,
				)
			} else {
				_, _ = b.client.ChannelMessageSend(
					reaction.ChannelID,
					"They've all been kicked!",
				)
			}
		} else {
			_, _ = b.client.ChannelMessageSend(
				reaction.ChannelID,
				"Report discarded.",
			)
			b.client.ChannelMessageDelete(
				reaction.ChannelID,
				reaction.MessageID,
			)
		}
		delete(b.reports, report.ReportID)
	}
}
