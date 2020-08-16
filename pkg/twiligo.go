package twiligo

import (
	"net/http"
	"path"
	"time"
)

const (
	baseURL             string = "https://api.twilio.com/2010-04-01"
	chatBaseURL         string = "https://chat.twilio.com/v2"
	conversationBaseURL string = "https://conversations.twilio.com/v1"
	proxyBaseURL        string = "https://proxy.twilio.com/v1"
)

// Exception represents an exception / bad response from the Twilio REST API.
type Exception struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Code     int    `json:"code"`
	MoreInfo string `json:"more_info"`
}

// Twilio holds the necessary important information for connecting to the Twilio REST API.
type Twilio struct {
	AccountSID string
	AuthToken  string
	HTTPClient *http.Client
}

// Error will print the current exception as a string.
func (e Exception) Error() string {
	return e.Message
}

// New creates a new instance of Twilio using the given Account SID and Auth Token.
func New(accountSID, authToken string) *Twilio {
	return NewCustomClient(accountSID, authToken, nil)
}

// NewCustomClient creates a new instance of Twilio using the given Account SID, Auth Token and HTTP Client.
func NewCustomClient(accountSID, authToken string, HTTPClient *http.Client) *Twilio {
	if HTTPClient == nil {
		HTTPClient = &http.Client{Timeout: time.Second * 30}
	}

	return &Twilio{
		AccountSID: accountSID,
		AuthToken:  authToken,
		HTTPClient: HTTPClient,
	}
}

func (twilio *Twilio) credentials() (string, string) {
	return twilio.AccountSID, twilio.AuthToken
}

func (twilio *Twilio) get(req *http.Request) (*http.Response, error) {
	return twilio.HTTPClient.Do(req)
}

func (twilio *Twilio) post(req *http.Request) (*http.Response, error) {
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	return twilio.HTTPClient.Do(req)
}

func (twilio *Twilio) delete(req *http.Request) (*http.Response, error) {
	return twilio.HTTPClient.Do(req)
}

func (twilio *Twilio) chatURL(resource string) string {
	return chatBaseURL + "/" + resource
}

func (twilio *Twilio) conversationURL(resource string) string {
	return conversationBaseURL + "/" + resource
}

func (twilio *Twilio) proxyURL(resource string) string {
	return proxyBaseURL + "/" + resource
}

func (twilio *Twilio) url(resource string) string {
	return baseURL + "/" + path.Join("Accounts", twilio.AccountSID, resource)
}
