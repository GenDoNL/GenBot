package Bot

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func (b *Bot) GetAccountCreationDate(user *discordgo.User) (timestamp int64) {
	id, _ := strconv.ParseUint(user.ID, 10, 64)
	timestamp = int64(((id >> 22) + 1420070400000) / 1000) // Divided by 1000 since we want seconds rather than ms
	return
}

func (b *Bot) ParseMention(tag string) (res string, err error) {
	idExp, err := regexp.Compile("<@!?([0-9]*)>")
	if err != nil {
		return
	}

	match := idExp.FindStringSubmatch(tag)
	if len(match) < 2 {
		err = errors.New("cannot parse id: " + tag)
		return
	}

	res = match[1]
	Log.Infof("Parsed %s to %s", tag, res)
	return
}

func (b *Bot) GetCommandTarget(s *discordgo.Session, m *discordgo.MessageCreate, data *ServerData, targetName string) (target *discordgo.User) {
	if len(m.Mentions) > 0 {
		target = m.Mentions[0]
		return
	}

	if targetName != "" {
		trg, err := b.getClosestUserByName(s, data, targetName)
		target = trg
		if err != nil {
			target = m.Author // Fallback if error occurs
		}

		// If a user is found through the state we will have to update that user.
		// Since State does not always contain the latest information on that user.
		updatedTarget, err := s.User(target.ID)
		if err != nil {
			return
		}
		target = updatedTarget
		return
	}

	target = m.Author
	return
}

// Helper functions regarding users
func (*Bot) getClosestUserByName(s *discordgo.Session, data *ServerData, user string) (foundUser *discordgo.User, err error) {
	currentMaxDistance := math.MaxInt64
	target := strings.ToLower(user)

	guild, err := s.Guild(data.ID)

	if err != nil {
		return
	}

	expensiveSubtitution := levenshtein.Options{
		InsCost: 1,
		DelCost: 1,
		SubCost: 3,
		Matches: levenshtein.IdenticalRunes,
	}

	for _, nick := range guild.Members {
		userName := strings.ToLower(nick.User.Username)

		levenDistance := levenshtein.DistanceForStrings([]rune(userName), []rune(target), expensiveSubtitution)

		if strings.Contains(userName, target) {
			levenDistance -= 10 // Bonus for starting with the correct name
		}

		nickDistance := math.MaxInt64
		// Prefer Server nickname over Discord username
		if nick.Nick != "" {
			userName = strings.ToLower(nick.Nick)
			nickDistance = levenshtein.DistanceForStrings([]rune(userName), []rune(target), expensiveSubtitution)

			if strings.Contains(userName, target) {
				nickDistance -= 10 // Bonus for starting with the correct name
			}
		}

		if levenDistance > nickDistance {
			levenDistance = nickDistance
		}

		if levenDistance < currentMaxDistance {
			currentMaxDistance = levenDistance
			foundUser = nick.User
		}
	}

	return
}