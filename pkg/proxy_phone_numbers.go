package twiligo

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

// AddPhoneNumberToProxyServiceOptions ...
type AddPhoneNumberToProxyServiceOptions struct {
	SID         string `url:"Sid,omitempty"`
	PhoneNumber string `url:",omitempty"`
	IsReserved  *bool  `url:",omitempty"`
}

// ProxyPhoneNumber ...
type ProxyPhoneNumber struct {
	SID          string                       `json:"sid"`
	AccountSID   string                       `json:"account_sid"`
	ServiceSID   string                       `json:"service_sid"`
	DateCreated  time.Time                    `json:"date_created"`
	DateUpdated  time.Time                    `json:"date_updated"`
	PhoneNumber  string                       `json:"phone_number"`
	FriendlyName string                       `json:"friendly_name"`
	IsoCountry   string                       `json:"iso_country"`
	Capabilities ProxyPhoneNumberCapabilities `json:"capabilities"`
	URL          string                       `json:"url"`
	IsReserved   bool                         `json:"is_reserved"`
	InUse        int                          `json:"in_use"`
}

// ProxyPhoneNumberCapabilities ...
type ProxyPhoneNumberCapabilities struct {
	MMSInbound    bool `json:"mms_inbound"`
	MMSOutbound   bool `json:"mms_outbound"`
	SMSInbound    bool `json:"sms_inbound"`
	SMSOutbound   bool `json:"sms_outbound"`
	VoiceInbound  bool `json:"voice_inbound"`
	VoiceOutbound bool `json:"voice_outbound"`
}

// AddPhoneNumberToProxyService ...
func (twilio *Twilio) AddPhoneNumberToProxyService(service string, options AddPhoneNumberToProxyServiceOptions) (*ProxyPhoneNumber, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	resource := "Services/" + service + "/PhoneNumbers"

	req, err := http.NewRequest(http.MethodPost, twilio.proxyURL(resource), strings.NewReader(params.Encode()))

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

	response := new(ProxyPhoneNumber)

	decoder.Decode(&response)

	return response, nil
}
