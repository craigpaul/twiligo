package twiligo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// ConversationOptions are all of the options that can be provided to a CreateNewConversation call.
type ConversationOptions struct {
	FriendlyName        string    `url:",omitempty"`
	DateCreated         time.Time `url:",omitempty"`
	DateUpdated         time.Time `url:",omitempty"`
	MessagingServiceSID string    `url:"MessagingServiceSid,omitempty"`
	Attributes          string    `url:",omitempty"`
	State               string    `url:",omitempty"`
	InactiveTimer       string    `url:"Timers.Inactive,omitempty"`
	ClosedTimer         string    `url:"Timers.Closed,omitempty"`
}

// Conversation represents a Twilio conversation between two or more connected participants.
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

// CreateNewConversation creates a new Conversation in Twilio with the provided options.
func (twilio *Twilio) CreateNewConversation(options ConversationOptions) (*Conversation, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	res, err := twilio.post(twilio.conversationURL("Conversations"), params)

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

// UpdateConversation will update an existing conversation in Twilio based on the provided identifier and options.
func (twilio *Twilio) UpdateConversation(conversationSID string, options ConversationOptions) (*Conversation, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	res, err := twilio.post(twilio.conversationURL("Conversations/"+conversationSID), params)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != http.StatusOK {
		err = new(Exception)

		decoder.Decode(err)

		return nil, err
	}

	response := new(Conversation)

	decoder.Decode(&response)

	return response, nil
}

// DeleteConversation will completely remove the conversation matching the given identifier from within Twilio.
func (twilio *Twilio) DeleteConversation(conversationSID string) error {
	res, err := twilio.delete(twilio.conversationURL("Conversations/" + conversationSID))

	if err != nil {
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		decoder := json.NewDecoder(res.Body)

		err = new(Exception)

		decoder.Decode(err)

		return err
	}

	return nil
}
