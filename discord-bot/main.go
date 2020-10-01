package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rivo/uniseg"

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

	// dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		if string(err.Error()) == "websocket: close 4004: Authentication failed." {
			fmt.Println("Make sure your Token is valid!")
			fmt.Println("Token:", Token)
		}
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

// from: https://github.com/bwmarrin/discordgo/wiki/FAQ#checking-if-a-message-is-a-direct-message-dm
func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

// tests if the user is the owner by comparing the userid to owner_id
func ownertest(s *discordgo.Session, m *discordgo.MessageCreate, owner_id string) {
	if m.Author.ID == owner_id {
		s.ChannelMessageSend(m.ChannelID, "you own the bot")
	} else {
		s.ChannelMessageSend(m.ChannelID, "you don't own the bot")
	}

}

func codesnippetprint(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	s.ChannelMessageSend(m.ChannelID, "```\n"+message+"\n```")
}

func dmtest(s *discordgo.Session, m *discordgo.MessageCreate) {
	dmResult, err := ComesFromDM(s, m)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "something went wrong")
	}
	if dmResult {
		s.ChannelMessageSend(m.ChannelID, "this is a dm")
	} else {
		s.ChannelMessageSend(m.ChannelID, "this isn't a dm")
	}

}

// runs neofetch
func neofetch(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmdstring := "/usr/bin/neofetch --stdout --color_blocks off | sed 's/\x1B[[0-9;?]*[a-zA-Z]//g' | sed '/^[[:space:]]*$/d'"

	fmt.Println("running neofetch")
	cmd := exec.Command("sudo", "su", "discord", "bash", "-c", cmdstring)

	out, err := cmd.CombinedOutput()
	codesnippetprint(s, m, string(out))
	if err != nil {
		error_str := string(err.Error())
		fmt.Println(error_str)
		codesnippetprint(s, m, error_str)
	}
	fmt.Println("ran neofetch")
}

// does uptime -p
func uptime(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmdstring := "uptime -p | sed '/^[[:space:]]*$/d'"

	fmt.Println("running uptime")
	cmd := exec.Command("bash", "-c", cmdstring)
	out, err := cmd.CombinedOutput()

	codesnippetprint(s, m, string(out))
	if err != nil {
		error_str := string(err.Error())
		fmt.Println(error_str)
		codesnippetprint(s, m, error_str)
	}

}

// kills the bot
func stopbot(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("shutting down bot")
	s.ChannelMessageSend(m.ChannelID, "shutting down bot!")
	os.Exit(0)
}

func temps(s *discordgo.Session, m *discordgo.MessageCreate) {
	// runs the temps command (my own shell script)
	cmdstring := "temps"

	fmt.Println("running temps")
	cmd := exec.Command("bash", "-c", cmdstring)
	out, err := cmd.CombinedOutput()

	codesnippetprint(s, m, string(out))
	if err != nil {
		error_str := string(err.Error())
		fmt.Println(error_str)
		codesnippetprint(s, m, error_str)
	}
}

// bash stuff, bc why not?
func bashRun(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmdstring := strings.Replace(m.Content, "!bash ", "", -1)

	cmd := exec.Command("sudo", "su", "discord", "bash", "-c", cmdstring)

	fmt.Println("running bash command:", cmdstring)
	out, err := cmd.CombinedOutput()

	if uniseg.GraphemeClusterCount(string(out)) > 2000 {
		s.ChannelMessageSend(m.ChannelID, "Sorry, the output of that command is over the 2000 character limit set by discord")
		return
	}

	codesnippetprint(s, m, string(out))
	if err != nil {
		error_str := string(err.Error())
		fmt.Println(error_str)
		codesnippetprint(s, m, error_str)
	}

	fmt.Println("finished running command!")
}

func noPermsCmd(s *discordgo.Session, m *discordgo.MessageCreate, owner_id string) {
	switch m.Content {
	case "!ping":
		s.ChannelMessageSend(m.ChannelID, "pong")

	case "!userid":
		// prints the id of the user doing !userid
		s.ChannelMessageSend(m.ChannelID, m.Author.ID)

	case "!ownertest":
		ownertest(s, m, owner_id)

	case "!dmtest":
		dmtest(s, m)

	case "!jebaited":
		// you just got jebaited
		s.ChannelMessageSend(m.ChannelID, "http://www.gardling.com/coolvideo6")

	case "!whyiuselinux":
		s.ChannelMessageSend(m.ChannelID, "http://www.gardling.com/whyiuselinux")

	case "bruh":
		// responds to every message "bruh" with "moment"
		s.ChannelMessageSend(m.ChannelID, "moment")

	case "!website":
		// links to my website
		s.ChannelMessageSend(m.ChannelID, "<http://www.gardling.com>")

	case "thx bot":
		// takes complements
		s.ChannelMessageSend(m.ChannelID, "np bro")

	case "pog":
		// responds to "pog" with "poggers"
		s.ChannelMessageSend(m.ChannelID, "poggers")

	case "!github":
		// links to the github of the bot
		s.ChannelMessageSend(m.ChannelID, "<https://github.com/Titaniumtown/golangstuff/tree/master/discord-bot>")

	case "!crab":
		// crab rave
		s.ChannelMessageSend(m.ChannelID, ":crab: :crab: :crab: :crab: :crab: :crab: :crab:")
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=LDU_Txk06tM")

	case "yeet":
		//test
		s.ChannelMessageSend(m.ChannelID, "itayayita")
		
	default:
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// if message sender is the bot, ignore message
	if m.Author.ID == s.State.User.ID {
		return
	}

	// my userid: 321028131982934017
	var owner_id = "321028131982934017"

	if m.Author.ID == owner_id {
		dmResult, err := ComesFromDM(s, m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "something went wrong")
		}
		if dmResult {
			if strings.HasPrefix(m.Content, "!bash") {
				bashRun(s, m)
			}

			switch m.Content {
			case "!neofetch":
				neofetch(s, m)

			case "!uptime":
				uptime(s, m)

			case "!stop":
				stopbot(s, m)

			case "!temps":
				temps(s, m)
			default:
				noPermsCmd(s, m, owner_id)
			}
		} else {
			noPermsCmd(s, m, owner_id)
		}
	} else {
		noPermsCmd(s, m, owner_id)
	}

	// send message: s.ChannelMessageSend(m.ChannelID, message)
	// delete message: s.ChannelMessageDelete(m.ChannelID, m.ID)

}
