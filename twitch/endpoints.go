package twitch

import (
	"log"
	"net/http"
	"net/url"
)

const baseURL string = "https://api.twitch.tv/helix"
const usersEndpoint string = "/users"
const webhooksSubscriptionsEndpoint string = "/webhooks/subscriptions"

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

func getWebhooksSubscriptions(page *string) (result *http.Request, err error) {
	result, err = http.NewRequest("GET", baseURL+webhooksSubscriptionsEndpoint, nil)
	if err != nil {
		log.Printf("Error creating getWebhooksSubscriptions request: %s", err)
		return
	}

	query := url.Values{}
	query.Add("first", "100")
	if page != nil {
		query.Add("after", *page)
	}
	result.URL.RawQuery = query.Encode()

	return
}
