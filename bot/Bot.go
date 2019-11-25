package main

import (
	"flag"
	_ "github.com/lib/pq"

	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/koffeinsource/go-imgur"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"time"
	"io/ioutil"
	"encoding/json"
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
	WebsiteUrl    string `json:"websiteurl"`
	OwmToken      string `json:"openweathermaptoken"`
	Mongo         string `json:"mongo"`
	Database      string `json:"database"`
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
	HServer    *HttpServer

	Modules map[string]Module

	AlbumCache map[string]*imgur.AlbumInfo
	imgClient  *imgur.Client

	mongoDB          *mongo.Client
	serverCollection *mongo.Collection
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

type Module interface {
	setup()
	execute(*discordgo.Session, *discordgo.MessageCreate, SentMessageData, *ServerData)
	retrieveCommands() map[string]Command
	retrieveHelp() (string, string)
}

func init() {
	flag.StringVar(&ConfigPath, "c", "./json/config.json", "<FILE>, where the configuration of the bot is located.")
	flag.Parse()
}

func main() {
	readConfig()
	startLogger()
	setupModules()
	startDataBase()
	initBot()
	return
}

func setupModules() {
	Modules = map[string]Module{}

	Modules["MetaModule"] = &MetaModule{}
	Modules["CommandModule"] = &CommandModule{}
	Modules["ModerationModule"] = &ModerationModule{}
	Modules["ImageModule"] = &ImageModule{}
	Modules["AnimeModule"] = &AnimeModule{}
	// Modules["OsuModule"] = &OsuModule{}, // Temporarily disable osu module until it is configurable

	for _, module := range Modules {
		module.setup()
	}
}

// Reads the config file.
func readConfig() {
	raw, err := ioutil.ReadFile(ConfigPath)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		json.Unmarshal(raw, &BotConfig)
		log.Info("Config file loaded.")
	}
}

func startLogger() {
	backEnd := logging.NewLogBackend(os.Stderr, "", 0)
	backEndFormatter := logging.NewBackendFormatter(backEnd, format)

	logging.SetBackend(backEndFormatter)

	log.Info("Logger Initialized")
}

func startDataBase() {
	var err error

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoDB, err = mongo.Connect(ctx, options.Client().ApplyURI(BotConfig.Mongo))

	if err != nil {
		log.Fatal(err.Error())
	}

	serverCollection = mongoDB.Database(BotConfig.Database).Collection("servers")
}

func initBot() {
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

	go HServer.start()

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

// EventHandler for when a message is sent.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID || len(m.Content) < 1 {
		return
	}

	msg := parseMessage(m)

	serverData := getServerDataFromChannel(s, m.ChannelID)

	if blocked, ok := serverData.BlockedCommands[msg.Command]; ok && blocked {
		return
	}

	for _, module := range Modules {
		module.execute(s, m, msg, serverData)
	}
}
