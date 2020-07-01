package internal

import (
	dg "github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
	"io/ioutil"
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
		log.Printf("Failed to authentication \"%s\", because\n%s", message.Author.ID, err.Error())
	}

	switch args[1] {
	case "add":
		if auth.IsAdmin {
			response := b.addRole(args[2], message.Author.ID, message.GuildID)
			_, _ = b.reply(message, response)
		} else {
			_, _ = b.reply(message, "You don't have permissions to run this command.")
		}
		break
	case "delete":
		if auth.IsAdmin {
			response := b.removeRole(args[2], message.Author.ID)
			_, _ = b.reply(message, response)
		} else {
			_, _ = b.reply(message, "You don't have permissions to run this command.")
		}
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
		return "Role doesn't exist"
	}

	if hasRole(roleID, b.config.Roles) {
		return "role was already in the config!"
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

	return "role added!"
}

func (b *Bot) removeRole(roleID string, userID string) string {
	if !hasRole(roleID, b.config.Roles) {
		return "role wasn't found and therefor couldn't be removed"
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

	return "role removed!"
}
