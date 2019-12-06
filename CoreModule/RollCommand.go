package CoreModule

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gendonl/genbot/Bot"
	"math/rand"
	"strconv"
	"strings"
	"time"
)


func initRollCommand() (cc CoreCommand) {
	cc = CoreCommand{
		name:        "roll",
		description: "Rolls a random number between 0 and the provided upper bound (0 and the upper bound included). Upper bound defaults to 100.",
		usage:       "`%sroll [upper-bound, inclusive]`",
		aliases:	 []string{},
		permission:  discordgo.PermissionSendMessages,
		execute:     (*CoreModule).rollCommand,
	}
	return
}

func (c *CoreModule) rollCommand(cmd CoreCommand, s *discordgo.Session, m *discordgo.MessageCreate, data *Bot.ServerData) {
	var maxRoll int64 = 100

	input := strings.SplitN(m.Content, " ", 3)

	if len(input) > 1 {
		newMaxRoll, err := strconv.ParseInt(input[1], 10, 64)
		if err == nil && newMaxRoll > 0 {
			maxRoll = newMaxRoll
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())

	result := strconv.FormatInt(rand.Int63n(maxRoll+1), 10)
	_, _ = s.ChannelMessageSend(m.ChannelID, result)
}