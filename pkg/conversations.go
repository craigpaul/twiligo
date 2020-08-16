package twiligo

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// CreateNewConversationOptions ...
type CreateNewConversationOptions struct {
	FriendlyName        string    `url:",omitempty"`
	DateCreated         time.Time `url:",omitempty"`
	DateUpdated         time.Time `url:",omitempty"`
	MessagingServiceSID string    `url:"MessagingServiceSid,omitempty"`
	Attributes          string    `url:",omitempty"`
	State               string    `url:",omitempty"`
	InactiveTimer       string    `url:"Timers.Inactive,omitempty"`
	ClosedTimer         string    `url:"Timers.Closed,omitempty"`
}

// Conversation ...
type Conversation struct {
	SID                 string    `json:"sid"`
	AccountSID          string    `json:"account_sid"`
	ChatServiceSID      string    `json:"chat_service_sid"`
	MessagingServiceSID string    `json:"messaging_service_sid"`
	FriendlyName        *string   `json:"friendly_name"`
	Attributes          string    `json:"attributes"`
	DateCreated         time.Time `json:"date_created"`
	DateUpdated         time.Time `json:"date_updated"`
	State               string    `json:"state"`
	Timers              struct {
		DateInactive string `json:"date_inactive"`
		DateClosed   string `json:"date_closed"`
	} `json:"timers"`
	Links struct {
		Participants string `json:"participants"`
		Messages     string `json:"messages"`
		Webhooks     string `json:"webhooks"`
	} `json:"links"`
	URL string `json:"url"`
}

// CreateNewConversation ...
func (twilio *Twilio) CreateNewConversation(options CreateNewConversationOptions) (*Conversation, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, twilio.conversationURL("Conversations"), strings.NewReader(params.Encode()))

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

	response := new(Conversation)

	decoder.Decode(&response)

	return response, nil
}