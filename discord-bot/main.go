package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rivo/uniseg"
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
		fmt.Println("User is an owner")
	} else {
		s.ChannelMessageSend(m.ChannelID, "you don't own the bot")
		fmt.Println("Use isn't an owner")
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
	//cmd := exec.Command("sudo", "su", "discord", "bash", "-c", cmdstring)
	cmd := exec.Command("bash", "-c", cmdstring)
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
	fmt.Println("ran uptime")

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
	fmt.Println("ran temps")
}

func ltsmstock(s *discordgo.Session, m *discordgo.MessageCreate) {
	//cmdstring := strings.Replace(m.Content, "!bash ", "", -1)
	var cmdstring = "python /mnt/hdd/ltsm_stocks/discordbot.py > /tmp/discordbot_ltsm.txt"
	cmd := exec.Command("bash", "-c", cmdstring)

	fmt.Println("running bash command:", cmdstring)
	out, err := cmd.CombinedOutput()

	//if uniseg.GraphemeClusterCount(string(out)) > 2000 {
	//	s.ChannelMessageSend(m.ChannelID, "Sorry, the output of that command is over the 2000 character limit set by discord")
	//	fmt.Println("Message sent was over 2000 char limit")
	//	return
	//}
	cmdstring = "cat /tmp/discordbot_ltsm.txt | grep filename"
	cmd = exec.Command("bash", "-c", cmdstring)

	fmt.Println("running bash command:", cmdstring)
	out, err = cmd.CombinedOutput()

	codesnippetprint(s, m, string(out))
	if err != nil {
		error_str := string(err.Error())
		fmt.Println(error_str)
		codesnippetprint(s, m, error_str)
	}

	fmt.Println("finished running command!")
}

// bash stuff, bc why not?
func bashRun(s *discordgo.Session, m *discordgo.MessageCreate) {
	cmdstring := strings.Replace(m.Content, "!bash ", "", -1)
	cmd := exec.Command("bash", "-c", cmdstring)

	fmt.Println("running bash command:", cmdstring)
	out, err := cmd.CombinedOutput()

	if uniseg.GraphemeClusterCount(string(out)) > 2000 {
		s.ChannelMessageSend(m.ChannelID, "Sorry, the output of that command is over the 2000 character limit set by discord")
		fmt.Println("Message sent was over 2000 char limit")
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
		fmt.Println("running ping")
		s.ChannelMessageSend(m.ChannelID, "pong")

	case "!userid":
		fmt.Println("running userid")
		// prints the id of the user doing !userid
		s.ChannelMessageSend(m.ChannelID, m.Author.ID)

	case "!ownertest":
		fmt.Println("Running ownertest")
		ownertest(s, m, owner_id)

	case "!dmtest":
		fmt.Println("running dmtest")
		dmtest(s, m)

	case "!jebaited":
		fmt.Println("running jebaited")
		// you just got jebaited
		s.ChannelMessageSend(m.ChannelID, "http://www.gardling.com/coolvideo6")

	case "!whyiuselinux":
		fmt.Println("running whyiuselinux")
		s.ChannelMessageSend(m.ChannelID, "http://www.gardling.com/whyiuselinux")

	case "bruh":
		fmt.Println("responding to 'bruh' with 'moment'")
		// responds to every message "bruh" with "moment"
		s.ChannelMessageSend(m.ChannelID, "moment")

	case "!website":
		fmt.Println("Linking to my website")
		// links to my website
		s.ChannelMessageSend(m.ChannelID, "<http://www.gardling.com>")

	case "thx bot":
		fmt.Println("responding to 'thx bot' with 'np bro'")
		// takes complements
		s.ChannelMessageSend(m.ChannelID, "np bro")

	case "pog":
		fmt.Println("responding to 'pog' with 'poggers'")
		// responds to "pog" with "poggers"
		s.ChannelMessageSend(m.ChannelID, "poggers")

	case "!github":
		fmt.Println("Linking to the sourcecode of the bot")
		// links to the github of the bot
		s.ChannelMessageSend(m.ChannelID, "<https://github.com/Titaniumtown/golangstuff/tree/master/discord-bot>")

	case "!crab":
		fmt.Println("crab rave time")
		// crab rave
		s.ChannelMessageSend(m.ChannelID, ":crab: :crab: :crab: :crab: :crab: :crab: :crab:")
		s.ChannelMessageSend(m.ChannelID, "https://www.youtube.com/watch?v=LDU_Txk06tM")

	case "yeet":
		fmt.Println("responding to 'yeet' with 'itayayita'")
		//test by a friend
		s.ChannelMessageSend(m.ChannelID, "itayayita")

	case "i cri":
		fmt.Println("responding to 'i cri' with 'are you shaking and crying rn?'")
		s.ChannelMessageSend(m.ChannelID, "are you shaking and crying rn?")

	default:
		return
	}
}

func printcreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	dt := time.Now()
	msginfo := fmt.Sprintf("(%s) server:%s channel:%s user:%s: %s\n", dt.String(), m.GuildID, m.ChannelID, m.Author.String(), m.Content)
	fmt.Println(msginfo)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// my userid: 321028131982934017
	var owner_id = "321028131982934017"
	switch m.Author.ID {
	case s.State.User.ID:
		return
	default:
		if rand.Intn(10000) == 69 {
			fmt.Sprintf("miceraft")
			s.ChannelMessageSend(m.ChannelID, "don't spam!")
		}
		printcreate(s, m)
		var githubPingGuildID = "795029627750973512"
		var githubPingChannelID = "795030212206264380"
		var githubPingChannelIDSendNormal = "795707486983815188"
		var githubPingChannelIDSendTest = "795697777706795018"
		var githubPingMessage = "<@&795688672418725908> new commits pushed to the master branch of TitaniumMC. :tada: Unless you have an extremely good reason not to update, You should really update your server ASAP! As always, you can download the latest build at: <http://www.gardling.com/titaniumclip.jar>"
		var githubPingChannelIDSend = githubPingChannelIDSendNormal

		githubwebhook := m.GuildID == githubPingGuildID && m.Author.String() == "GitHub#0000" && m.ChannelID == githubPingChannelID
		githubnotiftest := m.Content == "!githubnotificationtest"
		if githubwebhook || githubnotiftest {
			EmbedsString := fmt.Sprintf("%s", m.Embeds)
			fmt.Sprintf("# Github webhook Embed contents: (%s)", EmbedsString)
			var ciSkip = strings.Contains(EmbedsString, "[CI-SKIP]")

			if (strings.Contains(EmbedsString, " new commit ") && strings.Contains(EmbedsString, "TitaniumMC:master")) && !ciSkip {
				if githubnotiftest {
					githubPingChannelIDSend = githubPingChannelIDSendTest
				}
				s.ChannelMessageSend(githubPingChannelIDSend, githubPingMessage)
			}

		}
		if strings.Contains(m.Content, "vm.tiktok.com") {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
		}
		if m.Author.ID == owner_id {
			if m.Author.ID == owner_id {
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

				case "!stock":
					ltsmstock(s, m)

				default:
					noPermsCmd(s, m, owner_id)
				}
			} else {
				noPermsCmd(s, m, owner_id)
			}
		} else {
			noPermsCmd(s, m, owner_id)
		}
	}

	// send message: s.ChannelMessageSend(m.ChannelID, message)
	// delete message: s.ChannelMessageDelete(m.ChannelID, m.ID)

}
