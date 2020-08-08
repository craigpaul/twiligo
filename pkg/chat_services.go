package twiligo

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// ChatService ...
type ChatService struct {
	AccountSID                   string    `json:"account_sid"`
	DateCreated                  time.Time `json:"date_created"`
	DateUpdated                  time.Time `json:"date_updated"`
	DefaultChannelCreatorRoleSID string    `json:"default_channel_creator_role_sid"`
	DefaultChannelRoleSID        string    `json:"default_channel_role_sid"`
	DefaultServiceRoleSID        string    `json:"default_service_role_sid"`
	FriendlyName                 string    `json:"friendly_name"`
	Limits                       struct {
		ChannelMembers int `json:"channel_members"`
		UserChannels   int `json:"user_channels"`
	} `json:"limits"`
	Links struct {
		Channels string `json:"channels"`
		Users    string `json:"users"`
		Roles    string `json:"roles"`
		Bindings string `json:"bindings"`
	} `json:"links"`
	Notifications struct {
		RemovedFromChannel struct {
			Enabled bool `json:"enabled"`
		} `json:"removed_from_channel"`
		LogEnabled     bool `json:"log_enabled"`
		AddedToChannel struct {
			Enabled bool `json:"enabled"`
		} `json:"added_to_channel"`
		NewMessage struct {
			Enabled bool `json:"enabled"`
		} `json:"new_message"`
		InvitedToChannel struct {
			Enabled bool `json:"enabled"`
		} `json:"invited_to_channel"`
	} `json:"notifications"`
	Media struct {
		SizeLimitMB        int    `json:"size_limit_mb"`
		CompabilityMessage string `json:"compatibility_message"`
	} `json:"media"`
	PostWebhookURL         *string   `json:"post_webhook_url"`
	PreWebhookURL          *string   `json:"pre_webhook_url"`
	PreWebhookRetryCount   int       `json:"pre_webhook_retry_count"`
	PostWebhookRetryCount  int       `json:"post_webhook_retry_count"`
	ReachabilityEnabled    bool      `json:"reachability_enabled"`
	ReadStatusEnabled      bool      `json:"read_status_enabled"`
	SID                    string    `json:"sid"`
	TypingIndicatorTimeout int       `json:"typing_indicator_timout"`
	URL                    string    `json:"url"`
	WebhookFilters         *[]string `json:"webhook_filters"`
	WebhookMethod          *string   `json:"webhook_method"`
}

type createNewChatServiceOptions struct {
	FriendlyName string
}

// CreateNewChatService ...
func (twilio *Twilio) CreateNewChatService(name string) (*ChatService, error) {
	params, err := query.Values(createNewChatServiceOptions{FriendlyName: name})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, twilio.chatURL("Services"), strings.NewReader(params.Encode()))

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

	response := new(ChatService)

	decoder.Decode(&response)

	return response, nil
}
