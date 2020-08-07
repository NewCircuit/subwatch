// Watching has all the functions that do the actual "watching"
package internal

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"strconv"
	"strings"
	"time"
)

// This is called when the bot is ready. It's responsible for starting the interval.
func (bot *Bot) onReady(_ *discordgo.Session, ready *discordgo.Ready) {
	log.Printf("Ready as %s#%s\n", ready.User.Username, ready.User.Discriminator)
	bot.reviewGuild()
	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				bot.reviewGuild()
			}
		}
	}()
}

// this reviews the guild and looks for the
func (bot *Bot) reviewGuild() {
	var memberIDs []string
	var members string
	bot.reviewMembers("", &members, &memberIDs)

	// if there are members to kick, then send a report
	if len(memberIDs) > 0 {
		bot.startReport(members, memberIDs)
	} else {
		log.Println("I just reviewed the guild, no reports.")
	}

	// update the notification channel that the bot is still scanning.
	cest, _ := time.LoadLocation("Europe/Amsterdam")
	hour, min, _ := time.Now().In(cest).Clock()
	channelTopic := fmt.Sprintf("Last Checked: %d:", hour)
	if min < 10 {
		channelTopic += fmt.Sprintf("0%d", min)
	} else {
		channelTopic += strconv.Itoa(min)
	}
	channelTopic += " CEST"

	_, _ = bot.client.ChannelEditComplex(
		bot.config.NotificationChannel,
		&discordgo.ChannelEdit{
			Topic: channelTopic,
		},
	)
}

func (bot *Bot) startReport(summary string, memberIDs []string) {
	result := fmt.Sprintf(
		"**SubWatch Report**\n"+
			"These members will be kicked in %d seconds. "+
			"React with downvote to cancel.\n",
		bot.config.Delay/time.Second,
	)
	result += summary

	var msg *discordgo.Message
	var err error

	if len(result) < 2000 {
		msg, err = bot.client.ChannelMessageSend(
			bot.config.NotificationChannel,
			result,
		)
	} else {
		msg, err = bot.client.ChannelFileSend(
			bot.config.NotificationChannel,
			"list.txt",
			strings.NewReader(result),
		)
	}

	log.Println(result)

	if err != nil {
		log.Printf(
			"Failed to send a report to \"%s\" because\n%s\n",
			bot.config.NotificationChannel,
			err.Error(),
		)
		return
	}

	report := Report{
		MemberIDs: memberIDs,
		ReportID:  msg.ID,
		Cancel:    make(chan bool),
	}

	bot.reports[msg.ID] = report

	err = bot.client.MessageReactionAdd(
		msg.ChannelID,
		msg.ID,
		fmt.Sprintf(":voting:%s", bot.config.DownVote),
	)

	if err != nil {
		log.Printf(
			"Failed to react to \"%s\" in \"%s\" because\n%s\n",
			msg.ID,
			msg.ChannelID,
			err.Error(),
		)
	} else {
		go func() {
			timer := time.NewTimer(time.Second * bot.config.Delay)

			select {
			case <-timer.C:
				report.Cancel <- false
				break
			}
		}()
		toCancel := <-report.Cancel

		if !toCancel {
			result := bot.kickMembers(report.MemberIDs)

			if len(result) < 2000 {
				msg, err = bot.client.ChannelMessageSend(
					bot.config.NotificationChannel,
					result,
				)
			} else {
				msg, err = bot.client.ChannelFileSend(
					bot.config.NotificationChannel,
					"result.txt",
					strings.NewReader(result),
				)
			}

			if err != nil {
				log.Println("Failed to send conclusion to notification channel", err)
			}

		} else {
			_, _ = bot.client.ChannelMessageSend(
				bot.config.NotificationChannel,
				"Cancelled.",
			)
			delete(bot.reports, report.ReportID)
		}
	}
}

// This iterates through all the guild members.
func (bot *Bot) reviewMembers(memberID string, result *string, memberIDs *[]string) {
	members, err := bot.client.GuildMembers(bot.config.Guild, memberID, 1000)
	var lastMemberID string

	if err != nil {
		log.Printf(
			"Failed to fetch members for \"%s\" because\n%s\n",
			bot.config.Guild,
			err.Error(),
		)
		return
	}

	for _, member := range members {
		if !bot.checkRoles(member.Roles) {
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
		bot.reviewMembers(lastMemberID, result, memberIDs)
	}
}
