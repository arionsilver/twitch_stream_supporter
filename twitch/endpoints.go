package twitch

import (
	"log"
	"net/http"
	"net/url"
)

const baseURL string = "https://api.twitch.tv/helix"
const usersEndpoint string = "/users"

func getUsers(login string) (result *http.Request, err error) {
	result, err = http.NewRequest("GET", baseURL+usersEndpoint, nil)
	if err != nil {
		log.Printf("Error creating getUsers request: %s", err)
		return
	}

	query := url.Values{}
	query.Add("login", login)
	result.URL.RawQuery = query.Encode()

	return
}
