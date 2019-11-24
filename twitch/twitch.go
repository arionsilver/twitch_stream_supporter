package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Config config
type Config struct {
	CallbackServer string   `json:"callback_server"`
	Channels       []string `json:"channels"`
	Port           string   `json:"port"`
}

// Session session
type Session struct {
	authToken    string
	clientID     string
	clientSecret string
	client       http.Client
	config       Config

	tokenFile string
	token     tokenInfo
}

// UserInfo user info
type UserInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

// SubscriptionInfo subscription info
type SubscriptionInfo struct {
	Topic    string    `json:"topic"`
	Callback string    `json:"callback"`
	Expires  time.Time `json:"expires_at"`
}

type dataResult struct {
	Total      int               `json:"total"`
	Data       []json.RawMessage `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// NewSession creates session
func NewSession(auth, clientID, clientSecret string, config Config) (result Session) {
	result.authToken = auth
	result.clientID = clientID
	result.clientSecret = clientSecret
	result.client.Timeout = 10 * time.Second
	result.config = config

	return
}

func (session Session) addBearerAuth(req *http.Request) {
	req.Header = http.Header{}
	req.Header.Add("Authorization", "Bearer "+session.authToken)
}

func (session Session) addTokenAuth(req *http.Request) {
	req.Header = http.Header{}
	req.Header.Add("Authorization", "Bearer "+session.token.Token)
}

// GetUserInfo returns user info
func (session Session) GetUserInfo(userName string) (result UserInfo, err error) {
	req, err := getUsers(userName)
	if err != nil {
		return
	}

	session.addBearerAuth(req)
	res, err := session.client.Do(req)
	if err != nil {
		log.Printf("Error while execute GetUserInfo request: %s", err)
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("Request returned: %d", res.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return
	}

	var d dataResult
	if err = json.Unmarshal(body, &d); err != nil {
		log.Printf("Error while parsing response body: %s", err)
		log.Printf("Body: %s", body)
		return
	}

	for _, inner := range d.Data {
		// return first user, even if there are multiple results.
		err = json.Unmarshal(inner, &result)
		return
	}

	err = fmt.Errorf("User not found")
	return
}

// GetWebhooks get all webhooks subscriptions
func (session Session) GetWebhooks() (result []SubscriptionInfo, err error) {
	req, err := getWebhooksSubscriptions(nil)
	if err != nil {
		return
	}

	session.addTokenAuth(req)
	res, err := session.client.Do(req)
	if err != nil {
		log.Printf("Error while executing GetWebhooks request: %s", err)
		return
	}

	if res.StatusCode != 200 {
		err = fmt.Errorf("Request returned: %d", res.StatusCode)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error reading response body: %s", err)
		return
	}

	var d dataResult
	if err = json.Unmarshal(body, &d); err != nil {
		log.Printf("Error while parsing response body: %s", err)
		log.Printf("Body: %s", body)
		return
	}

	for _, inner := range d.Data {
		var sub SubscriptionInfo
		if err = json.Unmarshal(inner, &sub); err != nil {
			log.Printf("Error while parsing internal JSON: %s", err)
			log.Printf("Invalid JSON: %s", inner)
			err = nil
			continue
		}
		result = append(result, sub)
	}

	return
}

// SubscribeStream subscribe to stream by user ID
func (session Session) SubscribeStream(userID string, lease int) (err error) {
	type HubRequest struct {
		Callback string `json:"hub.callback"`
		Mode     string `json:"hub.mode"`
		Topic    string `json:"hub.topic"`
		Lease    int    `json:"hub.lease_seconds"`
	}
	body := HubRequest{}
	body.Callback = session.config.CallbackServer + session.config.Port
	body.Mode = "subscribe"
	body.Topic = subscribeStreamTopic(userID)
	body.Lease = lease

	bodyJSON, err := json.Marshal(body)
	if err != nil {
		log.Printf("Couldn't marshal the request JSON")
		return
	}

	req, err := postWebhookSubscription(bodyJSON)
	if err != nil {
		return
	}

	session.addTokenAuth(req)
	req.Header.Set("Content-Type", "application/json")
	res, err := session.client.Do(req)
	if err != nil {
		log.Printf("Error while executing SubscribeStream request: %s", err)
		return
	}

	bodyRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error while reading SubscribeStream request's body: %s", err)
		log.Printf("Body: %s", bodyRes)
		return
	}

	return nil
}
