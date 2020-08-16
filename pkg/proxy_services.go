package twiligo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// This constant is used to represent the type of number selection behaviour that the ProxyService should utilize.
const (
	AvoidSticky NumberSelectionBehaviour = iota + 1
	PreferSticky
)

// This constant is used to represent the matching level for where a proxy number must be located relative to a given participant.
const (
	AreaCode GeoMatchLevel = iota + 1
	Country
	ExtendedAreaCode
)

// CreateNewProxyServiceOptions are all of the options that can be provided to a CreateNewProxyService call.
type CreateNewProxyServiceOptions struct {
	CallbackURL              string                   `url:"CallbackUrl,omitempty"`
	ChatInstanceSID          string                   `url:",omitempty"`
	DefaultTTL               int                      `url:"DefaultTtl,omitempty"`
	GeoMatchLevel            GeoMatchLevel            `url:",omitempty"`
	InterceptCallbackURL     string                   `url:"InterceptCallbackUrl,omitempty"`
	NumberSelectionBehaviour NumberSelectionBehaviour `url:",omitempty"`
	OutOfSessionCallbackURL  string                   `url:"OutOfSessionCallbackUrl,omitempty"`
}

// GeoMatchLevel is used to define what matching level a ProxyService should implement for any given phone number belonging to the ProxyService.
type GeoMatchLevel int

// NumberSelectionBehaviour is used to define what number selection behaviour a ProxyService should implement for any given phone number belonging to the ProxyService.
type NumberSelectionBehaviour int

// ProxyService represents a Twilio Proxy Service that owns one or more proxy phone numbers, sessions, etc.
type ProxyService struct {
	SID                      string    `json:"sid"`
	AccountSID               string    `json:"account_sid"`
	ChatInstanceSID          *string   `json:"chat_instance_sid"`
	UniqueName               string    `json:"unique_name"`
	DefaultTTL               int       `json:"default_ttl"`
	CallbackURL              *string   `json:"callback_url"`
	GeoMatchLevel            string    `json:"geo_match_level"`
	NumberSelectionBehaviour string    `json:"number_selection_behaviour"`
	InterceptCallbackURL     *string   `json:"intercept_callback_url"`
	OutOfSessionCallbackURL  *string   `json:"out_of_session_callback_url"`
	DateCreated              time.Time `json:"date_created"`
	DateUpdated              time.Time `json:"date_updated"`
	URL                      string    `json:"url"`
	Links                    struct {
		Sessions     string `json:"sessions"`
		PhoneNumbers string `json:"phone_numbers"`
		ShortCodes   string `json:"short_codes"`
	} `json:"links"`
}

// CreateNewProxyService creates a new proxy service in Twilio.
func (twilio *Twilio) CreateNewProxyService(name string, options CreateNewProxyServiceOptions) (*ProxyService, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	params.Add("UniqueName", name)

	res, err := twilio.post(twilio.proxyURL("Services"), params)

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

	response := new(ProxyService)

	decoder.Decode(&response)

	return response, nil
}

func (geo GeoMatchLevel) String() string {
	return map[GeoMatchLevel]string{
		AreaCode:         "area-code",
		Country:          "country",
		ExtendedAreaCode: "extended-area-code",
	}[geo]
}

func (number NumberSelectionBehaviour) String() string {
	return map[NumberSelectionBehaviour]string{
		AvoidSticky:  "avoid-sticky",
		PreferSticky: "prefer-sticky",
	}[number]
}
