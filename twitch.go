package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "strings"

    "github.com/arionsilver/twitch_stream_supporter/twitch"
)

func startTwitchHelper(auth AuthInfo, tokenFile string, config twitch.Config, q chan bool) (c chan bool) {
    c = make(chan bool)
    go executeTwitchHelper(auth, tokenFile, config, c, q)

    return
}

func executeTwitchHelper(auth AuthInfo, tokenFile string, config twitch.Config, c chan bool, q chan bool) {
    defer func() { c <- true }()

    session := twitch.NewSession(auth.Twitch, auth.TwitchClientID, auth.TwitchClientSecret, config)
    validToken, err := session.CheckToken(tokenFile)
    if err != nil {
    	log.Printf("Error while validating twitch token. Quitting...\n%s", err)
    	return
    }

    if !validToken {
    	err = session.GenerateToken()
    	if err != nil {
    		log.Printf("Error while generating new twitch token. Quitting...\n%s", err)
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
    	subs, err := session.GetWebhooks()
    	if err != nil {
    		log.Printf("Error while fetching webhooks: %s", err)
    		return
    	}

    	for _, sub := range subs {
    		fmt.Printf("Subscription Information: %+v\n", sub)
    	}
    case "subscribe":
    	if len(args) > 2 {
    		id := args[1]
    		var lease int
    		if args[2] == "max" {
    			lease = 864000
    		} else {
    			lease, _ = strconv.Atoi(args[2])
    		}

    		if err := session.SubscribeStream(id, lease); err != nil {
    			log.Printf("Stream subscription failed: %s", err)
    		}
    	}
    }
}
