package Bot

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

// ServerData is the data which is saved for every server the bot is in.
type ServerData struct {
	ID              string                  `json:"id" bson:"id"`
	Commanders      map[string]bool         `json:"commanders" bson:"commanders"`
	Channels        map[string]*ChannelData `json:"channels" bson:"channels"`
	CustomCommands  map[string]*CommandData `json:"commands" bson:"commands"`
	BlockedCommands map[string]bool         `json:"blockedcommands" bson:"blockedcommands"`
	MeIrlData       map[string]*MeIrlData   `json:"me_irl" bson:"me_irl"`
	Key             string                  `json:"Key" bson:"Key"`
}

// CommandData is the data which is saved for every command
type CommandData struct {
	Name    string `json:"name" bson:"name"`
	Content string `json:"content" bson:"content"`
}

// ChannelData is the data which is saved for every channel
type ChannelData struct {
	ID     string   `json:"id" bson:"id"`
	Albums []string `json:"albums" bson:"albums"`
}

// MeIrlData is the data which is saved for every meIrlCommand
type MeIrlData struct {
	UserID   string `json:"id" bson:"id"`
	Nickname string `json:"nickname" bson:"nickname"`
	Content  string `json:"content" bson:"content"`
}

func (b *Bot) ServerDataFromChannel(s *discordgo.Session, channelID string) *ServerData {
	channel, _ := s.Channel(channelID)

	guildID := channel.GuildID

	return b.ServerDataFromID(guildID)
}

func (b *Bot) ServerDataFromID(guildID string) *ServerData {
	filter := bson.M{"id": guildID}
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	var data ServerData
	err := ServerCollection.FindOne(ctx, filter).Decode(&data)

	if err != nil {
		Log.Error(err)
		// If not found, create new server.
		data = newServerData(guildID)
		data.WriteToDB()
	} else {
		// Backwards compatibility, ensure no field is nil.
		if validateData(&data) {
			data.WriteToDB()
		}
	}

	return &data
}

func newServerData(guildID string) ServerData {
	data := ServerData{
		ID:              guildID,
		Key:             "!",
		CustomCommands:  map[string]*CommandData{},
		BlockedCommands: map[string]bool{},
		Channels:        map[string]*ChannelData{},
		Commanders:      map[string]bool{},
		MeIrlData:       map[string]*MeIrlData{},
	}
	return data
}

func validateData(data *ServerData) (updateRequired bool){
	updateRequired = false // Sanity check, since default value is false
	if data.CustomCommands == nil {
		data.CustomCommands = map[string]*CommandData{}
		updateRequired = true
	}

	if data.BlockedCommands == nil {
		data.BlockedCommands = map[string]bool{}
		updateRequired = true
	}

	if data.Channels == nil {
		data.Channels = map[string]*ChannelData{}
		updateRequired = true
	}

	if data.Commanders == nil {
		data.Commanders = map[string]bool{}
		updateRequired = true
	}

	if data.MeIrlData == nil {
		data.MeIrlData =  map[string]*MeIrlData{}
		updateRequired = true
	}
	return
}

func (data *ServerData) WriteToDB() (err error) {
	dataJson, err := bson.Marshal(data)

	if err != nil {
		Log.Error(err)
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	res, err := ServerCollection.ReplaceOne(
		ctx,
		bson.M{"id": data.ID},
		dataJson,
		options.Replace().SetUpsert(true))

	if err != nil {
		Log.Error(err)
		return err
	}

	Log.Infof("Updated %d server entries in database.", res.ModifiedCount)

	if res.UpsertedCount != 0 {
		Log.Infof("Inserted new server %s into database.", data.ID)
	}
	return
}

func updateData(serverId, index string, command interface{}) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// remove $ from input to clean it
	indexClean := strings.Replace(index, "$", "_", -1)

	_, err = ServerCollection.UpdateOne(
		ctx,
		bson.D{
			{"id", serverId },
		},
		bson.D{
				{"$set", bson.M{
					indexClean: &command,
				}},
		},
		options.Update().SetUpsert(true),
	)
	return
}

func (data *ServerData)  CreateCustomCommand(commandName, content string) (err error) {
	name := strings.ToLower(commandName)
	command := CommandData{name, content}
	commandClean := strings.Replace(commandName, ".", "_", -1)
	index := fmt.Sprintf("commands.%s", commandClean)
	return updateData(data.ID, index, command)
}

func (data *ServerData)  CreateMeIrl(meIrl MeIrlData) (err error) {
	commandClean := strings.Replace(meIrl.UserID, ".", "_", -1)
	index := fmt.Sprintf("me_irl.%s", commandClean)
	return updateData(data.ID, index, meIrl)
}

func (data *ServerData)  UpdateCommander(userId string, isCommander bool) (err error) {
	commandClean := strings.Replace(userId, ".", "_", -1)
	index := fmt.Sprintf("commanders.%s", commandClean)
	return updateData(data.ID, index, isCommander)
}

func (data *ServerData)  UpdateBlockedCommand(commandName string, isBlocked bool) (err error) {
	commandClean := strings.Replace(commandName, ".", "_", -1)
	index := fmt.Sprintf("blockedcommands.%s", commandClean)
	return updateData(data.ID, index, isBlocked)
}

func (data *ServerData) EditKey( key string) (err error) {
	index := "Key"
	return updateData(data.ID, index, key)
}

func removeData(serverId, index, name string) (deleted bool, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	// remove $ from input to clean it
	indexClean := strings.Replace(index, "$", "_", -1)

	result, err := ServerCollection.UpdateOne(
		ctx,
		bson.D{
			{"id", serverId },
		},
		bson.D{
			{"$unset", bson.M{
				indexClean: name,
			}},
		})
	if err != nil {
		return
	}
	if result.ModifiedCount > 0 {
		deleted = true
	}
	return
}

func (data *ServerData)  DeleteCustomCommand(commandName string) (deleted bool, err error) {
	commandClean := strings.Replace(commandName, ".", "_", -1)
	index := fmt.Sprintf("commands.%s", commandClean)
	return removeData(data.ID, index, commandName)
}

func (data *ServerData)  DeleteMeIrl(target string) (deleted bool, err error) {
	commandClean := strings.Replace(target, ".", "_", -1)
	index := fmt.Sprintf("me_irl.%s", commandClean)
	return removeData(data.ID, index, target)
}