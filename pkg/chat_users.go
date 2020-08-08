package twiligo

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// ChatUser ...
type ChatUser struct {
	SID                 string    `json:"sid"`
	AccountSID          string    `json:"account_sid"`
	ServiceSID          string    `json:"service_sid"`
	RoleSID             string    `json:"role_sid"`
	Identity            string    `json:"identity"`
	Attributes          string    `json:"attributes"`
	IsOnline            *bool     `json:"is_online"`
	IsNotifiable        *bool     `json:"is_notifiable"`
	FriendlyName        *string   `json:"friendly_name"`
	JoinedChannelsCount int       `json:"joined_channels_count"`
	DateCreated         time.Time `json:"date_created"`
	DateUpdated         time.Time `json:"date_updated"`
	Links               struct {
		UserChannels string `json:"user_channels"`
		USerBindings string `json:"user_bindings"`
	} `json:"links"`
	URL string `json:"url"`
}

// CreateNewChatUserOptions ...
type CreateNewChatUserOptions struct {
	RoleSID      string `url:"RoleSid,omitempty"`
	Attributes   string `url:",omitempty"`
	FriendlyName string `url:",omitempty"`
}

// CreateNewChatUser ...
func (twilio *Twilio) CreateNewChatUser(identity, serviceSID string, options CreateNewChatUserOptions) (*ChatUser, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	params.Add("Identity", identity)

	resource := "Services/" + serviceSID + "/Users"

	req, err := http.NewRequest(http.MethodPost, twilio.chatURL(resource), strings.NewReader(params.Encode()))

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(twilio.credentials())

	res, err := twilio.post(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusCreated {
		err = new(Exception)

		decoder.Decode(err)

		return nil, err
	}

	response := new(ChatUser)

	decoder.Decode(&response)

	return response, nil
}
