package twitch

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"
)

type tokenInfo struct {
    Token   string    `json:"token"`
    Expires time.Time `json:"expires"`
}

// CheckToken checks if token is valid
func (session *Session) CheckToken(tokenFile string) (bool, error) {
    session.tokenFile = tokenFile

    if _, err := os.Stat(tokenFile); err != nil {
    	log.Printf("Token file doesn't exist. Creating...")
    	err = ioutil.WriteFile(tokenFile, []byte{}, 0766)
    	return false, err
    }

    file, err := ioutil.ReadFile(tokenFile)
    if err != nil {
    	log.Printf("Error while reading token file: %s", err)
    	return false, err
    }

    info := tokenInfo{}
    if err = json.Unmarshal(file, &info); err != nil {
    	log.Printf("Error while parsing token file: %s", err)
    	return false, err
    }

    return session.validateToken(info)
}

// GenerateToken generates a Twitch app token
func (session *Session) GenerateToken() (err error) {
    req, err := postAppToken(session.clientID, session.clientSecret)
    if err != nil {
    	log.Printf("Fatal error creating request for app token generation: %s", err)
    	return
    }

    res, err := session.client.Do(req)
    if err != nil {
    	log.Printf("Fatal error requesting app token generation: %s", err)
    	return
    }

    if res.StatusCode != http.StatusOK {
    	err = fmt.Errorf("Response status code: %d", res.StatusCode)
    	return
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
    	log.Printf("Fatal error while reading response body: %s", err)
    	return
    }

    type tokenResponse struct {
    	Token        string   `json:"access_token"`
    	RefreshToken string   `json:"refresh_token"`
    	ExpiresIn    int      `json:"expires_in"`
    	Scopes       []string `json:"scope"`
    	TokenType    string   `json:"token_type"`
    }
    var response tokenResponse
    if err = json.Unmarshal(body, &response); err != nil {
    	log.Printf("Error while parsing response body: %s", err)
    	return
    }

    session.token.Expires = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
    session.token.Token = response.Token

    log.Println("App token generated.")
    log.Printf("\tExpires in: %s\n", session.token.Expires)
    log.Printf("\tRefresh Token: %s\n", response.RefreshToken)
    log.Printf("\tScopes: %+v\n", response.Scopes)
    log.Printf("\tToken Type: %s\n", response.TokenType)

    filebody, err := json.Marshal(session.token)
    if err != nil {
    	log.Printf("Error while marshalling token into json: %s", err)
    	return
    }

    err = ioutil.WriteFile(session.tokenFile, filebody, 0766)
    return
}

func (session *Session) validateToken(info tokenInfo) (bool, error) {
    if time.Now().After(info.Expires) {
    	return false, nil
    }

    req, err := validateAppToken(info.Token)
    if err != nil {
    	log.Printf("Failed in creating token validating request: %s", err)
    	return false, err
    }

    res, err := session.client.Do(req)
    if err != nil {
    	log.Printf("Token validation request failed: %s", err)
    	return false, err
    }

    if res.StatusCode != http.StatusOK {
    	err = fmt.Errorf("Token validation request failed with status code: %d", res.StatusCode)
    	return false, err
    }

    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
    	log.Printf("Fatal error while reading Token validation request's body: %s", err)
    	return false, err
    }

    type responseInfo struct {
    	ClientID  string `json:"client_id"`
    	Login     string `json:"login"`
    	UserID    string `json:"user_id"`
    	ExpiresIn int    `json:"expires_in"`
    }

    resultInfo := responseInfo{}
    if err = json.Unmarshal(body, &resultInfo); err != nil {
    	log.Printf("Fatal error while parsing response's json: %s", err)
    	log.Printf("Token validation body: %s", body)
    	return false, err
    }

    log.Print("Token validation successful")
    fmt.Printf("\tClient ID: %s\n", resultInfo.ClientID)
    fmt.Printf("\tExpires in(s): %d\n", resultInfo.ExpiresIn)
    fmt.Printf("\tExpiration date: %s\n", time.Now().Add(time.Duration(resultInfo.ExpiresIn)*time.Second))

    if resultInfo.ClientID != session.clientID {
    	err = fmt.Errorf("Token was validated but client ID is wrong: Considering this fatal")
    	return false, err
    }

    session.token = info

    return true, nil
}
