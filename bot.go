package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"strings"
	"time"
)

func (b *Bot) onRoleUpdate(session *discordgo.Session, member *discordgo.GuildMemberUpdate) {
	if !b.checkRoles(member.Roles) {
		b.sendEmbed(member.Member)
	}
}

func (b *Bot) onMessage(session *discordgo.Session, message *discordgo.MessageCreate) {
	if len(message.GuildID) == 0 || !strings.HasPrefix(message.Content, b.config.Prefix) {
		return
	}

	args := strings.Split(message.Content, " ")

	fmt.Println(args)

	if len(args) < 2 {
		return
	}

	switch args[1] {
	case "add":
		response := b.addRole(args[2], message.Author.ID, message.GuildID)
		b.client.ChannelMessageSend(message.ChannelID, response)
		break
	case "delete":
		response := b.removeRole(args[2], message.Author.ID)
		fmt.Println("delete")
		b.client.ChannelMessageSend(message.ChannelID, response)
		break
	}
}

func (b Bot) checkRoles(userRoles []string) bool {
	returnBool := false

	for _, item := range userRoles {
		if stringInSlice(item, b.config.Roles) {
			returnBool = true
		}
	}

	return returnBool
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func removeFromSlice(removeString string, list []string) []string {
	indexItem := 0

	for index, item := range list {
		if item == 	removeString {
			indexItem = index
		}
	}

	return append(list[:indexItem], list[indexItem+1:]...)
}

func (b Bot) sendEmbed(member *discordgo.Member) {
	embed := discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    member.User.Username,
			IconURL: member.User.AvatarURL(""),
		},
		Color: 0xff0000,
		Description: fmt.Sprintf("<@%s> needs to be checked up on", member.User.ID),
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     "Role watcher ⚠",
	}

	msg, _ := b.client.ChannelMessageSendEmbed(b.config.NotificationChannel, &embed)

	fmt.Println(msg.Content)
}

func (b *Bot) addRole(roleID string, userID string, guildID string) string {
	roleExists := false
	allRoles, _ := b.client.GuildRoles(guildID)

	for _, item := range allRoles {
		if item.ID == roleID {
			roleExists = true
			break
		}
	}

	if !roleExists {
		return fmt.Sprintf("<@%s> Role doesn't exist", userID)
	}

	if stringInSlice(roleID, b.config.Roles) {
		return fmt.Sprintf("<@%s> role was already in the config!", userID)
	}

	b.config.Roles = append(b.config.Roles, roleID)
	serialized, err := yaml.Marshal(b.config)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(configLocation, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("<@%s> role added!", userID)
}

func (b *Bot) removeRole(roleID string, userID string) string {
	if !stringInSlice(roleID, b.config.Roles) {
		return fmt.Sprintf("<@%s> role wasn't found and therefor couldn't be removed", userID)
	}

	b.config.Roles = removeFromSlice(roleID, b.config.Roles)

	serialized, err := yaml.Marshal(b.config)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(configLocation, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("<@%s> role removed!", userID)
}