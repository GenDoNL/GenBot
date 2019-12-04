package ModerationModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
)

func initLockCommand() (cc ModerationCommand) {
	cc = ModerationCommand{
		name:        "lock",
		description: "This command disallows the `everyone` role to speak in the current channel.",
		usage:       "`%slock`",
		aliases:	 []string{},
		permission:  discordgo.PermissionManageMessages,
		execute:     (*ModerationModule).lockCommand,
	}
	return
}

func initUnlockCommand() (cc ModerationCommand) {
	cc = ModerationCommand{
		name:        "unlock",
		description: "This command allows the `everyone` role to speak in the current channel.",
		usage:       "`%sunlock`",
		aliases:	 []string{},
		permission:  discordgo.PermissionManageMessages,
		execute:     (*ModerationModule).unlockCommand,
	}
	return
}

func (c *ModerationModule) lockCommand(cmd ModerationCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	ch, _ := s.Channel(m.ChannelID)
	sv, _ := s.Guild(data.ID)
	everyonePerms, err := c.Bot.GetRolePermissionsByName(ch, sv, "@everyone")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while trying to retrieve `everyone` role.")
		return
	}
	botPerm, _ := c.Bot.GetRolePermissions(c.Bot.BotID, ch.PermissionOverwrites)

	err = s.ChannelPermissionSet(ch.ID, c.Bot.BotID, "member", botPerm.Allow|0x800, botPerm.Deny&^0x800)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to lock channel, do I have the permissions to manage roles?")
		c.Bot.Log.Errorf("Error unlocking channel: %s", err)
		return
	}

	err = s.ChannelPermissionSet(ch.ID, everyonePerms.ID, everyonePerms.Type, everyonePerms.Allow&^0x800, everyonePerms.Deny|0x800)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to lock channel, do I have the permissions to manage roles?")
		c.Bot.Log.Errorf("Error unlocking channel: %s", err)
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

func (c *ModerationModule) unlockCommand(cmd ModerationCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	ch, _ := s.Channel(m.ChannelID)
	sv, _ := s.Guild(data.ID)
	everyonePerms, err := c.Bot.GetRolePermissionsByName(ch, sv, "@everyone")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Something went wrong while trying to retrieve `everyone` role.")
		return
	}

	err = s.ChannelPermissionSet(ch.ID, everyonePerms.ID, everyonePerms.Type, everyonePerms.Allow|0x800, everyonePerms.Deny&^0x800)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to unlock channel, do I have the permissions to manage roles?")
		c.Bot.Log.Errorf("Error unlocking channel: %s", err)
		return
	}

	s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
}

