package twitch

import (
    "bytes"
    "log"
    "net/http"
    "net/url"
)

const baseURL string = "https://api.twitch.tv/helix"
const streamsEndpoint string = "/streams"
const usersEndpoint string = "/users"
const webhooksSubscriptionsEndpoint string = "/webhooks/subscriptions"
const webhooksHubEndpoint string = "/webhooks/hub"

const baseOAuthURL string = "https://id.twitch.tv/oauth2"
const generateTokenEndpoint string = "/token"
const validateTokenEndpoint string = "/validate"

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

func postWebhookSubscription(body []byte) (result *http.Request, err error) {
    result, err = http.NewRequest("POST", baseURL+webhooksHubEndpoint, bytes.NewReader(body))
    if err != nil {
    	log.Printf("Error creating postWebhookSubscription request: %s", err)
    }

    return
}

func subscribeStreamTopic(userID string) string {
    query := url.Values{}
    query.Add("user_id", userID)
    return baseURL + streamsEndpoint + "?" + query.Encode()
}

func postAppToken(clientID, clientSecret string) (result *http.Request, err error) {
    result, err = http.NewRequest("POST", baseOAuthURL+generateTokenEndpoint, nil)
    if err != nil {
    	log.Printf("Error creating postAppToken request: %s", err)
    	return
    }

    query := url.Values{}
    query.Add("client_id", clientID)
    query.Add("client_secret", clientSecret)
    query.Add("grant_type", "client_credentials")
    result.URL.RawQuery = query.Encode()

    return
}

func validateAppToken(token string) (result *http.Request, err error) {
    result, err = http.NewRequest("GET", baseOAuthURL+validateTokenEndpoint, nil)
    if err != nil {
    	log.Printf("Error creating validateAppToken request: %s", err)
    	return
    }

    result.Header = http.Header{}
    result.Header.Add("Authorization", "OAuth "+token)

    return
}
