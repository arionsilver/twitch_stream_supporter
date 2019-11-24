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
	auth   string
	client http.Client
}

// UserInfo user info
type UserInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type dataResult struct {
	Data []json.RawMessage `json:"data"`
}

// NewSession creates session
func NewSession(auth string) (result Session) {
	result.auth = auth
	result.client.Timeout = 2 * time.Second

	return
}

func (session Session) addBearerAuth(req *http.Request) *http.Request {
	req.Header = http.Header{}
	req.Header.Add("Authorization", "Bearer "+session.auth)

	return req
}

// GetUserInfo returns user info
func (session Session) GetUserInfo(userName string) (result UserInfo, err error) {
	req, err := getUsers(userName)
	if err != nil {
		return
	}

	req = session.addBearerAuth(req)
	res, err := session.client.Do(req)
	if err != nil {
		log.Printf("Error while execute GetUserInfo request: %s", err)
	}

	if res.StatusCode != 200 {
		log.Printf("Request returned: %d", res.StatusCode)
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
