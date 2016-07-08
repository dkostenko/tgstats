package main

import (
	"errors"
	"flag"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jasonlvhit/gocron"
)

// Telegram bot token. Read more: https://core.telegram.org/bots
var (
	token       string
	channelName string
)

// Parse command line flags.
func init() {
	flag.StringVar(&token, "token", "", "Telegram bot token")
	flag.StringVar(&channelName, "channel", "", "Telegram channel name (with @)")

	flag.Parse()
}

func main() {
	err := validateCommandLineFlags()
	if err != nil {
		log.Println(err)
		return
	}

	collector, err := NewCollector(token, channelName, false)
	if err != nil {
		log.Println(err)
		return
	}

	// Collect information now.
	collector.GetChatMembersCount()

	// Collect information every 10 minutes.
	gocron.Every(10).Minutes().Do(collector.GetChatMembersCount)
	<-gocron.Start()
}

// Method validateCommandLineFlags return an error if it is.
func validateCommandLineFlags() error {
	if token == "" {
		return errors.New("Telegram bot token is required")
	}

	if channelName == "" {
		return errors.New("Telegram channel name is required")
	}

	return nil
}

// Collector get information about chanels with the help of the bot.
type Collector struct {
	Bot         *tgbotapi.BotAPI
	ChannelName string
}

// Collector constructor. Method returns new object of Collector struct.
// Input: Telegram bot token.
// Input: Telegram channel name.
// Input: Debug mode for Telegram bot connector.
func NewCollector(token, channelName string, debug bool) (*Collector, error) {
	var err error

	collector := &Collector{ChannelName: channelName}

	collector.Bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	collector.Bot.Debug = debug

	log.Println("Authorized on account", collector.Bot.Self.UserName)

	return collector, nil
}

// Method GetChatMembersCount print count of the channel members.
func (c *Collector) GetChatMembersCount() {
	count, err := c.Bot.GetChatMembersCount(tgbotapi.ChatConfig{SuperGroupUsername: c.ChannelName})
	if err != nil {
		log.Panic(err)
		return
	}

	log.Println("Chat members count:", count)
}
