package twiligo

import (
	"encoding/json"
	"net/http"

	dates "github.com/craigpaul/twiligo/internal"
	"github.com/google/go-querystring/query"
)

// This constant is used to represent the direction that a particular Message was sent.
const (
	Inbound MessageDirection = iota
	OutboundAPI
	OutboundCall
	OutboundReply
)

// This constant is used to represent the status of a particular Message.
const (
	Accepted MessageStatus = iota
	Queued
	Sending
	Sent
	Failed
	Delivered
	Undelivered
	Receiving
	Received
)

// CreateNewSMSMessageOptions are all of the options that can be provided to a CreateNewSMSMessage call.
type CreateNewSMSMessageOptions struct {
	Attempt             int    `url:"Attempt,omitempty"`
	From                string `url:"From,omitempty"`
	MessagingServiceSID string `url:"MessagingServiceSid,omitempty"`
	StatusCallback      string `url:"StatusCallback,omitempty"`
}

// Message represents any given type of message from Twilio.
type Message struct {
	SID                 string             `json:"sid"`
	AccountSID          string             `json:"account_sid"`
	APIVersion          string             `json:"api_version"`
	Body                string             `json:"body"`
	DateCreated         dates.Rfc2822Time  `json:"date_created"`
	DateSent            *dates.Rfc2822Time `json:"date_sent"`
	DateUpdated         dates.Rfc2822Time  `json:"date_updated"`
	Direction           MessageDirection   `json:"direction"`
	ErrorCode           *int               `json:"error_code"`
	ErrorDescription    *string            `json:"error_message"`
	From                string             `json:"from"`
	MessagingServiceSID *string            `json:"messaging_service_sid"`
	NumMedia            string             `json:"num_media"`
	NumSegments         string             `json:"num_segments"`
	Price               *string            `json:"price"`
	PriceUnit           *string            `json:"price_unit"`
	Status              MessageStatus      `json:"status"`
	SubresourceURIs     struct {
		Media string `json:"media"`
	} `json:"subresource_uris"`
	To  string `json:"to"`
	URI string `json:"uri"`
}

// MessageDirection is used to define the direction that a particular Message was sent.
type MessageDirection int

// MessageStatus is used to define the current status of a particular Message.
type MessageStatus int

// CreateNewSMSMessage sends an SMS message through Twilio using the given parameters.
func (twilio *Twilio) CreateNewSMSMessage(to, body string, options CreateNewSMSMessageOptions) (*Message, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	params.Add("To", to)
	params.Add("Body", body)

	res, err := twilio.post(twilio.url("Messages.json"), params)

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

	response := new(Message)

	decoder.Decode(&response)

	return response, nil
}

func (direction MessageDirection) String() string {
	return map[MessageDirection]string{
		Inbound:       "inbound",
		OutboundAPI:   "outbound-api",
		OutboundCall:  "outbound-call",
		OutboundReply: "outbound-reply",
	}[direction]
}

func (status MessageStatus) String() string {
	return map[MessageStatus]string{
		Accepted:    "accepted",
		Queued:      "queued",
		Sending:     "sending",
		Sent:        "sent",
		Failed:      "failed",
		Delivered:   "delivered",
		Undelivered: "undelivered",
		Receiving:   "receiving",
		Received:    "received",
	}[status]
}
