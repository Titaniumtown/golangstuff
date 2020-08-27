package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// var respond_self = true
	// if m.Content == "!respondself" {
	// 	if respond_self == false {
	// 		s.ChannelMessageSend(m.ChannelID, "ok, I'll respond to myself")
	// 		respond_self = true
	// 	} else if respond_self == true {
	// 		s.ChannelMessageSend(m.ChannelID, "ok, I won't respond to myself")
	// 		respond_self = false
	// 	}
	// }

	if m.Author.ID == s.State.User.ID {
		return
	}

	var sleep_time = time.Duration(rand.Int31n(1000))
	if m.Content == "ping" {
		time.Sleep(sleep_time * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

	if m.Content == "pong" {
		time.Sleep(sleep_time * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, "ping")
	}

}
