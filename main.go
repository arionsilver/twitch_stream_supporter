package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
)

// AuthInfo auth info
type AuthInfo struct {
	Discord            string `json:"discord"`
	Twitch             string `json:"twitch"`
	TwitchClientID     string `json:"twitch-client-id"`
	TwitchClientSecret string `json:"twitch-client-secret"`
}

func readAuthFile(filename string) (auth AuthInfo) {
	if filename == "" {
		log.Fatal("Authentication information filename must not be empty.")
	}

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("An error occured while reading authentication file: %s", err)
	}

	if err = json.Unmarshal(file, &auth); err != nil {
		log.Fatalf("An error occured while parsing authentication file: %s", err)
	}

	return
}

func main() {
	var authFile string

	flag.StringVar(&authFile, "auth", "", "authentication information file for both twitch and discord")
	flag.Parse()

	auth := readAuthFile(authFile)

	twitchQuitSignal := make(chan bool)
	discordQuitSignal := make(chan bool)

	twitchExit := startTwitchHelper(auth, twitchQuitSignal)
	discordExit := startDiscordBot(auth.Discord, discordQuitSignal)
	quitSignal := make(chan os.Signal)
	signal.Notify(quitSignal, os.Interrupt, os.Kill)

	select {
	case quit := <-twitchExit:
		if quit {
			log.Print("Twitch helper called exit")
			discordQuitSignal <- true
		}
	case quit := <-discordExit:
		if quit {
			log.Print("Discord helper called exit")
			twitchQuitSignal <- true
		}
	case <-quitSignal:
		twitchQuitSignal <- true
		discordQuitSignal <- true
		log.Print("Interrupt called")
	}
}
