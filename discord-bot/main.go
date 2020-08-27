package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func main() {
	var token_file = "token.txt"
	if !(fileExists(token_file)) {
		fmt.Println("token file:", token_file, "was not found, so the file was created, please put your bot token there")
		emptyFile, err := os.Create(token_file)
		if err != nil {
			log.Fatal(err)
		}
		emptyFile.Close()
		os.Exit(1)
	}
	content, err := ioutil.ReadFile(token_file)
	Token := string(content)

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
		if string(err.Error()) == "websocket: close 4004: Authentication failed." {
			fmt.Println("Make sure your Token is valid!")
			fmt.Println("Token:", Token)
		}
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// if message sender is the bot, ignore message
	if m.Author.ID == s.State.User.ID {
		return
	}

	// var sleep_time = time.Duration(rand.Int31n(1000))
	switch m.Content {
	case "ping":
		// time.Sleep(sleep_time * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, "pong")

	case "pong":
		// time.Sleep(sleep_time * time.Millisecond)
		s.ChannelMessageSend(m.ChannelID, "ping")
	}

	// command tests:
	// if strings.HasPrefix(m.content, "!")

}
