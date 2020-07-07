package internal

import (
	"github.com/go-yaml/yaml"
	"io/ioutil"
)

// add a new required role. the bot's response is what's returned
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
