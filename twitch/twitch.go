package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Session session
type Session struct {
	authToken    string
	clientID     string
	clientSecret string
	client       http.Client
}

// UserInfo user info
type UserInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type dataResult struct {
	Total      int               `json:"total"`
	Data       []json.RawMessage `json:"data"`
	Pagination struct {
		Cursor string `json:"cursor"`
	} `json:"pagination"`
}

// NewSession creates session
func NewSession(auth, clientID, clientSecret string) (result Session) {
	result.authToken = auth
	result.clientID = clientID
	result.clientSecret = clientSecret
	result.client.Timeout = 10 * time.Second

	return
}

func (session Session) addBearerAuth(req *http.Request) {
	req.Header = http.Header{}
	req.Header.Add("Authorization", "Bearer "+session.authToken)
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
func (session Session) GetWebhooks() (err error) {
	req, err := getWebhooksSubscriptions(nil)
	if err != nil {
		return
	}

	session.addBearerAuth(req)
	res, err := session.client.Do(req)
	if err != nil {
		log.Printf("Error while executing GetWebhooks request: %s", err)
		return
	}

	log.Printf("%+v", res)

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

	fmt.Printf("%+v\n", d)

	return
}
