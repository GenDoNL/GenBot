package main

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"math"
	"net/url"
	"strconv"
	"strings"
	"unicode"
	"context"
	"time"
	"go.mongodb.org/mongo-driver/bson"
)

// This function parse a discord.MessageCreate into a SentMessageData struct.
func parseMessage(m *discordgo.MessageCreate) SentMessageData {
	// Remove all white-space characters, except for new-lines.
	f := func(c rune) bool {
		return c != '\n' && unicode.IsSpace(c)
	}

	split := strings.FieldsFunc(m.Content, f)
	key := m.Content[:1]
	commandName := strings.ToLower(split[0][1:])
	if (len(commandName) > 1) && (commandName[len(commandName)-1] == '\n') {
		commandName = commandName[:len(commandName)-1]
	}

	content := split[1:]

	return SentMessageData{key, commandName, content, m.ID, m.ChannelID, m.Mentions, m.Author}
}

// This function a string into an ID if the string is a mention.
func parseMention(str string) (string, error) {
	if len(str) < 5 || (string(str[0]) != "<" || string(str[1]) != "@" || string(str[len(str)-1]) != ">") {
		return "", errors.New("error while parsing mention, this is not an user")
	}

	res := str[2 : len(str)-1]

	// Necessary to allow nicknames.
	if string(res[0]) == "!" {
		res = res[1:]
	}

	return res, nil
}

// Returns the ServerData of a server, given a message object.
func getServerDataDB(s *discordgo.Session, channelID string) *ServerData {
	channel, _ := s.Channel(channelID)

	servID := channel.GuildID

	//var err error
	filter := bson.M{"id": servID}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	var data ServerData
	err := serverCollection.FindOne(ctx, filter).Decode(&data)


	if err != nil {
		log.Error(err.Error())
		data = ServerData{ID: servID, Key: "!"}

		err2 := writeServerDataDB(&data)
		if err2 != nil {
			log.Fatal(err2)
		}
	}

	return &data

	//if len(Servers) == 0 {
	//	Servers = make(map[string]*ServerData)
	//}
	//
	//if serv, ok := Servers[servID]; ok {
	//	return serv
	//}

	//Servers[servID] = &ServerData{ID: servID, Key: "!"}
	//return Servers[servID]

}

func writeServerDataDB(data *ServerData) error {
	dataB, err := bson.Marshal(data)

	if err != nil {
		log.Error(err)
		return  err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res, err2 := serverCollection.ReplaceOne(ctx, bson.M{"id": data.ID}, dataB)

	log.Infof("Updated %d entries in database.", res)

	if res.MatchedCount == 0 {
		serverCollection.InsertOne(ctx, dataB)
		log.Infof("Inserted new server %s into database.", data.ID)
	}

	if err2 != nil {
		log.Error(err)
		return  err
	}

	return nil
}

// Checks whether a user id (String) is in a slice of users.
func userInSlice(a string, list []*discordgo.User) bool {
	for _, b := range list {
		if b.ID == a {
			return true
		}
	}
	return false
}

// Gets the specific role by name out of a role list.
func getRoleByName(name string, roles []*discordgo.Role) (r discordgo.Role, e error) {
	for _, elem := range roles {
		if elem.Name == name {
			r = *elem
			return
		}
	}
	e = errors.New("Role name not found in the specified role array: " + name)
	return
}

// Gets the permission override object from a role id.
func getRolePermissions(id string, perms []*discordgo.PermissionOverwrite) (p discordgo.PermissionOverwrite, e error) {
	for _, elem := range perms {
		if elem.ID == id {
			p = *elem
			return
		}
	}
	e = errors.New("permissions not found in the specified role: " + id)
	return
}

func getRolePermissionsByName(ch *discordgo.Channel, sv *discordgo.Guild, name string) (p discordgo.PermissionOverwrite, e error) {
	//get role object for given name
	role, _ := getRoleByName(name, sv.Roles)
	return getRolePermissions(role.ID, ch.PermissionOverwrites)
}

func getRoleById(s *discordgo.Session, data *ServerData, id string) (*discordgo.Role, error) {
	g, _ := s.Guild(data.ID)
	for _, role := range g.Roles {
		if role.ID == id {
			println(role.Name)
			return role, nil
		}
	}
	return nil, errors.New("role not found in list")
}

// isValidUrl tests a string to determine if it is a url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

// Creates a command in the given server given a name and a message.
func createCommand(data *ServerData, commandName, message string) error {
	name := strings.ToLower(commandName)
	if strings.Contains(name, "\n") {
		log.Info("Trying to add command name with newline, aborted.")
		return errors.New("trying to add command with a name that contains a new line")
	}
	data.CustomCommands[name] = &CommandData{name, message}
	writeServerDataDB(data)
	return nil
}

func checkCommandsMap(data *ServerData) {
	if len(data.CustomCommands) == 0 {
		data.CustomCommands = make(map[string]*CommandData)
	}
}

func checkChannelsMap(data *ServerData) {
	if len(data.Channels) == 0 {
		data.Channels = make(map[string]*ChannelData)
	}
}

func findLastMessageWithAttachOrEmbed(s *discordgo.Session, msg SentMessageData, amount int) (result string, e error) {
	msgList, _ := s.ChannelMessages(msg.ChannelID, amount, msg.MessageID, "", "")

	for _, x := range msgList {
		if len(x.Embeds) > 0 {
			result = x.Embeds[0].URL
			e = nil
			return
		} else if len(x.Attachments) > 0 {
			result = x.Attachments[0].URL
			e = nil
			return
		}
	}

	result = ""
	e = errors.New("Unable to find message with attachment or embed")
	return
}

func getAccountCreationDate(user *discordgo.User) (timestamp int64) {
	id, _ := strconv.ParseUint(user.ID, 10, 64)
	timestamp = int64(((id >> 22) + 1420070400000) / 1000) // Divided by 1000 since we want seconds rather than ms
	return
}

func getClosestUserByName(s *discordgo.Session, data *ServerData, user string) (foundUser *discordgo.User, err error) {
	currentMaxDistance := math.MaxInt64
	target := strings.ToLower(user)

	guild, err := s.Guild(data.ID)

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

func getCommandTarget(s *discordgo.Session, msg SentMessageData, data *ServerData) (target *discordgo.User) {
	if len(msg.Mentions) > 0 {
		target = msg.Mentions[0]
	} else {
		if len(msg.Content) > 0 {
			trg, err := getClosestUserByName(s, data, strings.Join(msg.Content, " "))
			target = trg
			if err != nil {
				target = msg.Author // Fallback if error occurs
			}
		} else {
			target = msg.Author
		}
	}
	return
}
