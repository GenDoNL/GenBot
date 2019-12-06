package Bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
)

func (b *Bot) GetRoleByName(name string, roles []*discordgo.Role) (r discordgo.Role, e error) {
	for _, elem := range roles {
		if elem.Name == name {
			r = *elem
			return
		}
	}
	e = errors.New("Role name not in slice " + name)
	return
}

// Gets the permission override object from a role id.
func (b *Bot) GetRolePermissions(id string, perms []*discordgo.PermissionOverwrite) (p discordgo.PermissionOverwrite, e error) {
	for _, elem := range perms {
		if elem.ID == id {
			p = *elem
			return
		}
	}
	e = errors.New("permissions of role not found " + id)
	return
}

func (b *Bot) GetRolePermissionsByName(ch *discordgo.Channel, sv *discordgo.Guild, name string) (p discordgo.PermissionOverwrite, e error) {
	//get role object for given name
	role, _ := b.GetRoleByName(name, sv.Roles)
	return b.GetRolePermissions(role.ID, ch.PermissionOverwrites)
}