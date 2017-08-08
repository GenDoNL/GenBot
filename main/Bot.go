package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"net/http"

	"github.com/koffeinsource/go-imgur"
	"github.com/koffeinsource/go-klogger"

	"github.com/thehowl/go-osuapi"
)

// Variables used for Command line parameters
var (
	Token        string
	BotID        string
	ImgurID      string
	OsuID        string
	DataLocation string

	Servers         map[string]*ServerData
	AdminCommands   map[string]func(*discordgo.Session, MessageData, *ServerData)
	PruneCommands   map[string]func(*discordgo.Session, MessageData, *ServerData)
	DefaultCommands map[string]func(*discordgo.Session, MessageData, *ServerData)
	AlbumCache      map[string]*imgur.AlbumInfo

	osuClient osuapi.Client
	imgClient *imgur.Client
)

// ServerData is the data which is saved for every server the bot is in.
type ServerData struct {
	ID            string                  `json:"id"`
	Commanders    map[string]bool         `json:"commanders"`
	Channels      map[string]*ChannelData `json:"channels"`
	Commands      map[string]*CommandData `json:"commands"`
	meIrlCommands map[string]*MeIrlData   `json:"meIrlCommand"`
	Key           string                  `json:"Key"`
}

// ChannelData is the data which is saved for every channel
type ChannelData struct {
	ID     string   `json:"ID"`
	Albums []string `json:"albums"`
}

// MeIrlData is the data which is saved for every meIrlCommand
type MeIrlData struct {
	UserID   string `json:"id"`
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
}

// CommandData is the data which is saved for every command
type CommandData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// MessageData is the parsed message, allows for easier access to arguments.
type MessageData struct {
	Key       string
	Command   string
	Content   []string
	MessageID string
	ChannelID string
	Mentions  []*discordgo.User
	Author    *discordgo.User
}

func setUp() {
	// Set up all the commands
	AdminCommands = map[string]func(*discordgo.Session, MessageData, *ServerData){
		"addcommand":   addCommandCommand,
		"delcommand":   delCommandCommand,
		"addcommander": addCommanderCommand,
		"delcommander": delCommanderCommand,
		"setkey":       setKeyCommand,
		"addalbum":     addAlbumCommand,
		"delalbum":     delAlbumCommand,
		"forcereload":  forceGetAlbumCommand,
		"addme_irl":    addMeIrlCommand,
		"delme_irl":    delMeIrlCommand,
		"lock":         lockChannelCommand,
		"unlock":       unlockChannelCommand,
	}

	PruneCommands = map[string]func(*discordgo.Session, MessageData, *ServerData){
		"prune": pruneCommand,
	}

	DefaultCommands = map[string]func(*discordgo.Session, MessageData, *ServerData){
		"i":           getImageCommand,
		"image":       getImageCommand,
		"meIrl":       meIrlCommand,
		"help":        helpCommand,
		"commandlist": commandListCommands,
	}

	// Set up the cache so we do not have to make multiple API calls for the same album.
	AlbumCache = make(map[string]*imgur.AlbumInfo)
}

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ImgurID, "i", "", "Imgur Token")
	flag.StringVar(&OsuID, "o", "", "Osu Token")
	flag.StringVar(&DataLocation, "d", "./data.json", "Data Location")
	flag.Parse()
}

func main() {
	imgClient = new(imgur.Client)
	imgClient.HTTPClient = new(http.Client)
	imgClient.Log = new(klogger.CLILogger)
	imgClient.ImgurClientID = ImgurID

	osuClient = *osuapi.NewClient(OsuID)

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
		return
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	readServerData()
	setUp()

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID || len(m.Content) < 1 {
		return
	}

	// Parse the command into a MessageData struct
	msg := parseCommand(m)

	// Retrieve the server data of the given server.
	serverData := getServerData(s, m.ChannelID)

	//Check for any beatmap links in message
	checkBeatmapLink(s, m)

	if serverData.Key != msg.Key {
		return
	}

	// Check if a Command is in one of the maps.
	if f, ok := AdminCommands[msg.Command]; ok {
		perm, _ := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if isCommander, ok := serverData.Commanders[m.Author.ID]; isAdmin(perm) || (ok && isCommander) {
			f(s, msg, serverData)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Insufficient permission to use **"+msg.Command+"**.")
		}
	} else if f, ok := PruneCommands[msg.Command]; ok {
		perm, _ := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
		if isCommander, ok := serverData.Commanders[m.Author.ID]; isAllowedToPrune(perm) || (ok && isCommander) {
			f(s, msg, serverData)
		} else {
			s.ChannelMessageSend(m.ChannelID, "Insufficient permission to use **"+msg.Command+"**.")
		}
	} else if f, ok := DefaultCommands[msg.Command]; ok {
		f(s, msg, serverData)
	} else if cmd, ok := serverData.Commands[msg.Command]; ok {
		s.ChannelMessageSend(m.ChannelID, cmd.Content)
	}
}
