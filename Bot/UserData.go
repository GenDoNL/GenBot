package Bot

import (
	"context"
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type UserData struct {
	UserId string `bson:"userid"`
	AniListData AniListUserData `bson:"anilistdata"`
}

type AniListUserData struct {
	UserId int `bson:"aniuserid"`
	LastUpdated time.Time `bson:"lastupdated"`
}

func newUserData(userID string) UserData {
	data := UserData{
		UserId: userID,
		AniListData: AniListUserData{},
	}
	return data
}

func UserDataFromMessage(m *discordgo.MessageCreate) (userData *UserData) {
	input := strings.SplitN(m.Content, " ", 2)
	if len(input) == 1 || len(m.Mentions) > 0 {
		// No user was mentioned, this means author is the target
		if len(input) == 1 {
			userData = UserDataFromID(m.Author.ID)
		} else { // User was mentioned, make this user the target.
			userData = UserDataFromID(m.Mentions[0].ID)
		}
	}
	return
}

func UserDataFromID(userID string) *UserData {
	filter := bson.M{"userid": userID}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var data UserData
	err := UserCollection.FindOne(ctx, filter).Decode(&data)

	if err != nil {
		Log.Error(err)
		data = newUserData(userID)
	}

	return &data
}

func (data *UserData) WriteToDB() (err error) {
	dataJson, err := bson.Marshal(data)

	if err != nil {
		Log.Error(dataJson)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res, err :=  UserCollection.ReplaceOne(
		ctx,
		bson.M{"userid": data.UserId},
		dataJson,
		options.Replace().SetUpsert(true))

	if err != nil {
		Log.Error(err)
		return err
	}

	Log.Infof("Updated %d user entries in database.", res.ModifiedCount)

	if res.UpsertedCount != 0 {
		Log.Infof("Inserted new user %s into database.", data.UserId)
	}

	return
}



//// discordgo.User helpers ////
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