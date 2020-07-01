package twiligo

import (
	"net/http"
	"path"
	"time"
)

const baseURL string = "https://api.twilio.com/2010-04-01"

// Exception ...
type Exception struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"more_info"`
}

// Twilio ...
type Twilio struct {
	AccountSid string
	AuthToken  string
	HTTPClient *http.Client
}

// Error ...
func (e Exception) Error() string {
	return e.Message
}

// New ...
func New(accountSid, authToken string) *Twilio {
	return NewCustomClient(accountSid, authToken, nil)
}

// NewCustomClient ...
func NewCustomClient(accountSid, authToken string, HTTPClient *http.Client) *Twilio {
	if HTTPClient == nil {
		HTTPClient = &http.Client{Timeout: time.Second * 30}
	}

	return &Twilio{
		AccountSid: accountSid,
		AuthToken:  authToken,
		HTTPClient: HTTPClient,
	}
}

func (twilio *Twilio) credentials() (string, string) {
	return twilio.AccountSid, twilio.AuthToken
}

func (twilio *Twilio) get(req *http.Request) (*http.Response, error) {
	return twilio.HTTPClient.Do(req)
}

func (twilio *Twilio) url(resource string) string {
	return baseURL + "/" + path.Join("Accounts", twilio.AccountSid, resource)
}
