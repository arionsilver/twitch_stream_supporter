package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"

    "github.com/arionsilver/twitch_stream_supporter/twitch"
    "github.com/bwmarrin/discordgo"
)

func startDiscordBot(auth string, config twitch.Config, q chan bool) (c chan bool) {
    c = make(chan bool)
    session, err := discordgo.New(auth)
    if err != nil {
    	log.Printf("Error while starting discord session: %s", err)
    	go func() { c <- true }()
    	return
    }
    go executeDiscordBot(session, config, c, q)
    go createSimpleHTTPServer(session, config)

    return
}

func executeDiscordBot(session *discordgo.Session, config twitch.Config, c chan bool, q chan bool) {
    defer func() { c <- true }()

    if err := session.Open(); err != nil {
    	log.Printf("Couldn't open discord session: %s", err)
    	return
    }

    defer session.Close()
    session.AddHandler(func(s *discordgo.Session, event *discordgo.GuildCreate) {
    	for _, channel := range event.Guild.Channels {
    		fmt.Printf("Channel ID= %s GuildID= %s Name= %s\n", channel.ID, channel.GuildID, channel.Name)
    	}
    })

    <-q // wait on quit
}

type rootHandler struct {
    session *discordgo.Session
    config  twitch.Config
}

func (handler rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer w.WriteHeader(200)
    r.ParseForm()
    if len(r.Form.Get("hub.challenge")) > 0 {
    	w.Write([]byte(r.Form.Get("hub.challenge")))
    } else if r.ContentLength > 0 {
    	body, err := ioutil.ReadAll(r.Body)
    	if err != nil {
    		log.Printf("Failed while reading request body: %s", err)
    		return
    	}

    	type jsonRequest struct {
    		Data []struct {
    			ID	string `json:"id"`
    			Title    string `json:"title"`
    			UserName string `json:"user_name"`
    		} `json:"data"`
    	}

    	log.Printf("Body: %s", body)

    	var info jsonRequest
    	if err = json.Unmarshal(body, &info); err != nil {
    		log.Printf("Couldn't parse the request as a JSON: %s", err)
    		return
    	}

    	for _, channel := range handler.config.Channels {
    		for _, data := range info.Data {
    			content := fmt.Sprintf("%s has started streaming.\nStream title: **%s**\nhttp://twitch.tv/%s", data.UserName, data.Title, data.UserName)
    			handler.session.ChannelMessageSend(channel, content)
    		}
    	}
    } else {
    	log.Printf("Empty request received.\n\tHeaders: %+v\n\tRemote address: %s", r.Header, r.RemoteAddr)
    }
}

func createSimpleHTTPServer(session *discordgo.Session, config twitch.Config) {
    http.Handle("/", rootHandler{session, config})
    log.Fatal(http.ListenAndServe(config.Port, nil))
}
