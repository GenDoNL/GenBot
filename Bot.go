package main

import (
	"flag"
	_ "github.com/lib/pq"

	"github.com/bwmarrin/discordgo"
	"github.com/koffeinsource/go-imgur"
	"github.com/op/go-logging"
	"os"
)

// Set up Logger
var log = logging.MustGetLogger("GenBot")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// The Config of the bot
type Config struct {
	BotToken      string `json:"bottoken"`
	OsuToken      string `json:"osutoken"`
	ImgurToken    string `json:"imgurtoken"`
	SauceNaoToken string `json:"saucenaotoken"`
	DataLocation  string `json:"datalocation"`
}

// MessageData is the parsed message, allows for easier access to arguments.
type SentMessageData struct {
	Key       string
	Command   string
	Content   []string
	MessageID string
	ChannelID string
	Mentions  []*discordgo.User
	Author    *discordgo.User
}

type Command struct {
	Name        string
	Description string
	Usage       string
	Permission  int
	Execute     func(Command, *discordgo.Session, SentMessageData, *ServerData)
}

// Variables used for Command line parameters
var (
	ConfigPath string
	BotID      string
	BotConfig  Config

	Servers   map[string]*ServerData
	Modules   []Module
	CmdModule CommandModule

	AlbumCache map[string]*imgur.AlbumInfo
	imgClient  *imgur.Client
)

// ServerData is the data which is saved for every server the bot is in.
type ServerData struct {
	ID             string                  `json:"id"`
	Commanders     map[string]bool         `json:"commanders"`
	Channels       map[string]*ChannelData `json:"channels"`
	CustomCommands map[string]*CommandData `json:"commands"`
	MeIrlData      map[string]*MeIrlData   `json:"me_irl"`
	Key            string                  `json:"Key"`
}

// CommandData is the data which is saved for every command
type CommandData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// ChannelData is the data which is saved for every channel
type ChannelData struct {
	ID     string   `json:"id"`
	Albums []string `json:"albums"`
}

// MeIrlData is the data which is saved for every meIrlCommand
type MeIrlData struct {
	UserID   string `json:"id"`
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
}

type Module interface {
	setup()
	execute(*discordgo.Session, *discordgo.MessageCreate)
}

func init() {
	flag.StringVar(&ConfigPath, "c", "./json/config.json", "<FILE>, where the configuration of the bot is located.")
	flag.Parse()
}

func main() {
	readConfig()
	startLogger()
	readServerData()
	setupModules()
	initializeBot()
	return
}

func setupModules() {
	Modules = []Module{
		&CmdModule,
		&OsuModule{},
	}

	for _, module := range Modules {
		module.setup()
	}
}

func startLogger() {
	backEnd := logging.NewLogBackend(os.Stderr, "", 0)
	backEndFormatter := logging.NewBackendFormatter(backEnd, format)

	logging.SetBackend(backEndFormatter)

	log.Info("Logger Initialized")
}

func initializeBot() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + BotConfig.BotToken)
	if err != nil {
		log.Critical("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		log.Critical("error obtaining account details,", err)
		return
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		log.Critical("error opening connection,", err)
		return
	}

	log.Infof("Bot %s is now running.  Press CTRL-C to exit.", BotID)
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

	log.Debug("Received message: " + m.Content)

	for _, module := range Modules {
		module.execute(s, m)
	}
}
