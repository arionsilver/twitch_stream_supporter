package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/arionsilver/twitch_stream_supporter/twitch"
)

func startTwitchHelper(auth AuthInfo, q chan bool) (c chan bool) {
	c = make(chan bool)
	go executeTwitchHelper(auth, c, q)

	return
}

func executeTwitchHelper(auth AuthInfo, c chan bool, q chan bool) {
	defer func() { c <- true }()

	session := twitch.NewSession(auth.Twitch, auth.TwitchClientID, auth.TwitchClientSecret)
	validToken, err := session.CheckToken()
	if err != nil {
		log.Printf("Error while validating twitch token. Quitting...")
		return
	}

	if !validToken {
		err = session.GenerateToken()
		if err != nil {
			log.Printf("Error while generating new twitch token. Quitting...")
			return
		}
	}

	go waitOnInput(session)
	<-q // wait on quit
}

func waitOnInput(session twitch.Session) {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		if err := stdin.Err(); err != nil {
			log.Printf("Error occurred reading stdin!?: %s", err)
		}

		text := stdin.Text()
		split := strings.Split(text, " ")
		if len(split) > 0 {
			executeCommand(session, split)
		}
	}
}

func executeCommand(session twitch.Session, args []string) {
	switch args[0] {
	case "user":
		if len(args) > 1 {
			userName := args[1]
			user, err := session.GetUserInfo(userName)
			if err != nil {
				log.Printf("Error while fetch user info: %s", err)
				return
			}

			fmt.Printf("User info:\n\tID: %s\n\tDisplay Name: %s\n", user.ID, user.DisplayName)
		}
	case "webhooks":
		err := session.GetWebhooks()
		if err != nil {
			log.Printf("Error while fetching webhooks: %s", err)
			return
		}
	}
}
