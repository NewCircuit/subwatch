package internal

import (
	"fmt"
	dg "github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"strings"
)

func (b *Bot) onMessage(_ *dg.Session, message *dg.MessageCreate) {
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
		_, _ = b.client.ChannelMessageSend(message.ChannelID, response)
		break
	case "delete":
		response := b.removeRole(args[2], message.Author.ID)
		fmt.Println("delete")
		_, _ = b.client.ChannelMessageSend(message.ChannelID, response)
		break
	}
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

	if hasRole(roleID, b.config.Roles) {
		return fmt.Sprintf("<@%s> role was already in the config!", userID)
	}

	b.config.Roles = append(b.config.Roles, roleID)
	serialized, err := yaml.Marshal(b.config)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(b.confLoc, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("<@%s> role added!", userID)
}

func (b *Bot) removeRole(roleID string, userID string) string {
	if !hasRole(roleID, b.config.Roles) {
		return fmt.Sprintf("<@%s> role wasn't found and therefor couldn't be removed", userID)
	}

	b.config.Roles = removeFromSlice(roleID, b.config.Roles)

	serialized, err := yaml.Marshal(b.config)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(b.confLoc, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("<@%s> role removed!", userID)
}
