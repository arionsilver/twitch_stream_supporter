package twitch

import "fmt"

// CheckToken checks if token is valid
func (session Session) CheckToken() (result bool, err error) {
	result = false
	err = fmt.Errorf("Not Implemented")

	return
}

// GenerateToken generates a Twitch app token
func (session Session) GenerateToken() (err error) {
	err = fmt.Errorf("Not Implemented")
	return
}
