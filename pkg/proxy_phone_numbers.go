package twiligo

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/go-querystring/query"
)

// AddPhoneNumberToProxyServiceOptions are all of the options that can be provided to a AddPhoneNumberToProxyService call.
type AddPhoneNumberToProxyServiceOptions struct {
	SID         string `url:"Sid,omitempty"`
	PhoneNumber string `url:",omitempty"`
	IsReserved  *bool  `url:",omitempty"`
}

// ProxyPhoneNumber represents an IncomingPhoneNumber that has been attached to a ProxyService within Twilio.
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

// ProxyPhoneNumberCapabilities describes the capabilities afforded to a given ProxyPhoneNumber according to Twilio.
type ProxyPhoneNumberCapabilities struct {
	MMSInbound    bool `json:"mms_inbound"`
	MMSOutbound   bool `json:"mms_outbound"`
	SMSInbound    bool `json:"sms_inbound"`
	SMSOutbound   bool `json:"sms_outbound"`
	VoiceInbound  bool `json:"voice_inbound"`
	VoiceOutbound bool `json:"voice_outbound"`
}

// AddPhoneNumberToProxyService attaches a phone number to the given proxy service in Twilio.
func (twilio *Twilio) AddPhoneNumberToProxyService(serviceSID string, options AddPhoneNumberToProxyServiceOptions) (*ProxyPhoneNumber, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	res, err := twilio.post(twilio.proxyURL("Services/"+serviceSID+"/PhoneNumbers"), params)

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

// RemovePhoneNumberFromProxyService remove the given IncomingPhoneNumber from the given ProxyService within Twilio.
func (twilio *Twilio) RemovePhoneNumberFromProxyService(serviceSID, phoneNumberSID string) error {
	res, err := twilio.delete(twilio.proxyURL("Services/" + serviceSID + "/PhoneNumbers/" + phoneNumberSID))

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
