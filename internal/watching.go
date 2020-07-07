// Watching has all the functions that do the actual "watching"
package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"time"
)

// This is called when the bot is ready. It's responsible for starting the interval.
func (b *Bot) onReady(_ *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Ready as %s#%s\n", ready.User.Username, ready.User.Discriminator)
	b.reviewGuild()
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				b.reviewGuild()
			}
		}
	}()
}

// this reviews the guild and looks for the
func (b *Bot) reviewGuild() {
	result := "**__SubWatch__** ⚠\n" +
		"These people need to be checked up on:\n" +
		" * upvote: **notify & kick**\n" +
		" * downvote: **discard**\n"
	members := ""
	memberIDs := []string{}
	b.reviewMembers("", &members, &memberIDs)

	if len(members) > 0 {
		msg, err := b.client.ChannelMessageSend(
			b.config.NotificationChannel,
			result+members,
		)

		if err != nil {
			log.Printf(
				"Failed to send a report to \"%s\" because\n%s\n",
				b.config.NotificationChannel,
				err.Error(),
			)
			return
		}

		b.reports[msg.ID] = Report{
			MemberIDs: memberIDs,
			ReportID:  msg.ID,
		}

		_ = b.client.MessageReactionAdd(
			msg.ChannelID,
			msg.ID,
			fmt.Sprintf(":voting:%s", b.config.UpVote),
		)
		err = b.client.MessageReactionAdd(
			msg.ChannelID,
			msg.ID,
			fmt.Sprintf(":voting:%s", b.config.DownVote),
		)

		if err != nil {
			log.Printf(
				"Failed to react to \"%s\" in \"%s\" because\n%s\n",
				msg.ID,
				msg.ChannelID,
				err.Error(),
			)
		}
	} else {
		b.client.ChannelMessageSend(
			b.config.NotificationChannel,
			"**__SubWatch__**\nNo one to check-up on.",
		)
	}
}

// This iterates through all the guild members.
func (b *Bot) reviewMembers(memberID string, result *string, memberIDs *[]string) {
	members, err := b.client.GuildMembers(b.config.Guild, memberID, 1000)
	var lastMemberID string

	if err != nil {
		log.Printf(
			"Failed to fetch members for \"%s\" because\n%s\n",
			b.config.Guild,
			err.Error(),
		)
		return
	}

	for _, member := range members {
		if !b.checkRoles(member.Roles) {
			*memberIDs = append(*memberIDs, member.User.ID)
			*result += fmt.Sprintf(
				" - %s#%s (<@%s>)\n",
				member.User.Username,
				member.User.Discriminator,
				member.User.ID,
			)
		}
		lastMemberID = member.User.ID
	}

	if len(members) == 1000 {
		b.reviewMembers(lastMemberID, result, memberIDs)
	}
}
