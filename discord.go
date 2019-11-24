package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func startDiscordBot(auth string, q chan bool) (c chan bool) {
	c = make(chan bool)
	go executeDiscordBot(auth, c, q)

	return
}

func executeDiscordBot(auth string, c chan bool, q chan bool) {
	defer func() { c <- true }()

	session, err := discordgo.New(auth)
	if err != nil {
		log.Printf("Error while starting discord session: %s", err)
		return
	}

	if err := session.Open(); err != nil {
		log.Printf("Couldn't open discord session: %s", err)
		return
	}

	defer session.Close()

	<-q // wait on quit
}
