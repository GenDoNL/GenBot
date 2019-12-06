package Bot

import (
	"context"
	"encoding/json"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/op/go-logging"
	"go.mongodb.org/mongo-driver/mongo"
)

type Module interface {
	Execute(*discordgo.Session, *discordgo.MessageCreate, *ServerData)
	HelpFields() (title string, content string)
	CommandInfo(name string, data *ServerData) (response *Embed, found bool)
}

type Server interface {
	Start()
	GetUrlFromID(guildID string) string
}

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

type Bot struct {
	Log        *logging.Logger
	ConfigPath string
	BotID      string
	Config     Config
	Modules    []Module
	Server	   Server

	// DB
	ServerCollection *mongo.Collection

}

func (b *Bot) initLogger() {
	b.Log = logging.MustGetLogger("GenBot")
	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	backEnd := logging.NewLogBackend(os.Stderr, "", 0)
	backEndFormatter := logging.NewBackendFormatter(backEnd, format)

	logging.SetBackend(backEndFormatter)

	b.Log.Info("Logger initialized.")
}

func (b *Bot) startDataBase() {
	var err error

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mongoDB, err := mongo.Connect(ctx, options.Client().ApplyURI(b.Config.Mongo))

	if err != nil {
		b.Log.Fatal(err.Error())
	}

	b.ServerCollection = mongoDB.Database(b.Config.Database).Collection("servers")
	b.Log.Info("Database initialized.")
}


func (b *Bot) InitBot(m []Module, site Server, configPath string) {
	b.Modules = m
	b.ConfigPath = configPath
	b.Server = site
	b.initLogger()
	b.readConfig()
	b.startDataBase()

	go b.Server.Start()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + b.Config.BotToken)
	if err != nil {
		b.Log.Critical("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		b.Log.Critical("error obtaining account details,", err)
		return
	}

	b.BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(b.messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		b.Log.Critical("error opening connection,", err)
		return
	}

	b.Log.Infof("Bot %s is now running.  Press CTRL-C to exit.", b.BotID)
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// EventHandler for when a message is sent.
func (b *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Whatever happens, the bot should not go down when a single faulty message is received.
	defer func() {
		if r := recover(); r != nil {
			b.Log.Criticalf("Bot panicked in channel %s on message: `%s`, %s", m.ChannelID, m.Content, r)
		}
	}()

	// Ignore all messages created by the bot itself or other bots.
	// Also ignore messages of length 0 for now.
	if m.Author.ID == b.BotID || m.Author.Bot ||  len(m.Content) < 1   {
		return
	}

	ch, err := s.Channel(m.ChannelID)

	if err != nil || ch.Type != 0 {
		s.ChannelMessageSend(ch.ID, "I do currently not work in DMs.")
		return
	}

	data := b.ServerDataFromChannel(s, m.ChannelID)

	for _, module := range b.Modules {
		module.Execute(s, m, data)
	}
}

// Reads the config file.
func (b *Bot) readConfig() {
	raw, err := ioutil.ReadFile(b.ConfigPath)
	if err != nil {
		b.Log.Fatal(err.Error())
	}

	err = json.Unmarshal(raw, &b.Config)
	if err != nil {
		b.Log.Fatal(err.Error())
	}

	b.Log.Info("Config file loaded.")
}
