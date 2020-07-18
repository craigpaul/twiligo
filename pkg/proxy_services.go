package twiligo

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// ...
const (
	AvoidSticky NumberSelectionBehaviour = iota + 1
	PreferSticky
)

// ...
const (
	AreaCode GeoMatchLevel = iota + 1
	Country
	ExtendedAreaCode
)

// CreateNewProxyServiceOptions ...
type CreateNewProxyServiceOptions struct {
	CallbackURL              string                   `url:"CallbackUrl,omitempty"`
	ChatInstanceSid          string                   `url:",omitempty"`
	DefaultTTL               int                      `url:"DefaultTtl,omitempty"`
	GeoMatchLevel            GeoMatchLevel            `url:",omitempty"`
	InterceptCallbackURL     string                   `url:"InterceptCallbackUrl,omitempty"`
	NumberSelectionBehaviour NumberSelectionBehaviour `url:",omitempty"`
	OutOfSessionCallbackURL  string                   `url:"OutOfSessionCallbackUrl,omitempty"`
}

// GeoMatchLevel ...
type GeoMatchLevel int

// NumberSelectionBehaviour ...
type NumberSelectionBehaviour int

// ProxyService ...
type ProxyService struct {
	Sid                      string    `json:"sid"`
	AccountSid               string    `json:"account_sid"`
	ChatInstanceSid          *string   `json:"chat_instance_sid"`
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

// CreateNewProxyService ...
func (twilio *Twilio) CreateNewProxyService(name string, options CreateNewProxyServiceOptions) (*ProxyService, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	params.Add("UniqueName", name)

	req, err := http.NewRequest(http.MethodPost, twilio.proxyURL("Services"), strings.NewReader(params.Encode()))

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
