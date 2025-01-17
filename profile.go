package messenger

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Profile struct holds data associated with Facebook profile
type Profile struct {
	Name           string  `json:"name,omitempty"`     // Instagram Fallback
	UserName       string  `json:"username,omitempty"` // Instagram username
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	ProfilePicture string  `json:"profile_pic,omitempty"`
	Locale         string  `json:"locale,omitempty"`
	Timezone       float64 `json:"timezone,omitempty"`
	Gender         string  `json:"gender,omitempty"`
}

// GetProfile fetches the recipient's profile from facebook platform
// Non empty UserID has to be specified in order to receive the information
func (m *Messenger) GetProfile(userID string, scope string) (*Profile, error) {
	url := GraphAPI + "/v14.0/" + userID

	if scope != "" {
		url += "?fields=" + scope
	}

	resp, err := m.doRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured: " + er.Error.Message)
	}
	profile := new(Profile)
	return profile, json.Unmarshal(read, profile)
}

type accountLinking struct {
	//Recipient is Page Scoped ID
	Recipient string `json:"recipient"`
}

// GetPSID fetches user's page scoped id during authentication flow
// one must supply a valid and not expired authentication token provided by facebook
// https://developers.facebook.com/docs/messenger-platform/account-linking/authentication
func (m *Messenger) GetPSID(token string) (*string, error) {
	resp, err := m.doRequest("GET", fmt.Sprintf(GraphAPI+"/v10.0/me?fields=recipient&account_linking_token=%s", token), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	read, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		er := new(rawError)
		json.Unmarshal(read, er)
		return nil, errors.New("Error occured: " + er.Error.Message)
	}
	acc := new(accountLinking)
	return &acc.Recipient, json.Unmarshal(read, acc)
}
