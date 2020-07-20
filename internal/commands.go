package internal

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// add a new required role. the bot's response is what's returned
func (bot *Bot) addRole(roleID string, guildID string) string {
	roleExists := false
	allRoles, _ := bot.client.GuildRoles(guildID)

	for _, item := range allRoles {
		if item.ID == roleID {
			roleExists = true
			break
		}
	}

	if !roleExists {
		return "Role doesn't exist"
	}

	if hasRole(roleID, bot.config.Roles) {
		return "role was already in the config!"
	}

	bot.config.Roles = append(bot.config.Roles, roleID)
	serialized, err := yaml.Marshal(bot.config)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(bot.confLoc, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return "role added!"
}

func (bot *Bot) removeRole(roleID string) string {
	if !hasRole(roleID, bot.config.Roles) {
		return "role wasn't found and therefor couldn't be removed"
	}

	bot.config.Roles = removeFromSlice(roleID, bot.config.Roles)

	serialized, err := yaml.Marshal(bot.config)

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(bot.confLoc, serialized, 0660)

	if err != nil {
		panic(err)
	}

	return "role removed!"
}
