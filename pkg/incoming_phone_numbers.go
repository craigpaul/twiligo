package twiligo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	dates "github.com/craigpaul/twiligo/internal"
	"github.com/google/go-querystring/query"
)

// CreateNewIncomingPhoneNumberOptions are all of the options that can be provided to a CreateNewIncomingPhoneNumber call.
type CreateNewIncomingPhoneNumberOptions struct {
	AddressSID           string `url:",omitempty"`
	AreaCode             string `url:",omitempty"`
	BundleSID            string `url:",omitempty"`
	EmergencyAddressSID  string `url:",omitempty"`
	EmergencyStatus      string `url:",omitempty"`
	FriendlyName         string `url:",omitempty"`
	IdentitySID          string `url:",omitempty"`
	PhoneNumber          string `url:",omitempty"`
	SmsApplicationSID    string `url:",omitempty"`
	SmsFallbackMethod    string `url:",omitempty"`
	SmsFallbackURL       string `url:"SmsFallbackUrl,omitempty"`
	SmsMethod            string `url:",omitempty"`
	SmsURL               string `url:"SmsUrl,omitempty"`
	StatusCallback       string `url:",omitempty"`
	StatusCallbackMethod string `url:",omitempty"`
	TrunkSID             string `url:",omitempty"`
	VoiceApplicationSID  string `url:",omitempty"`
	VoiceFallbackMethod  string `url:",omitempty"`
	VoiceFallbackURL     string `url:"VoiceFallbackUrl,omitempty"`
	VoiceMethod          string `url:"VoiceMethod,omitempty"`
	VoiceReceiveMode     string `url:",omitempty"`
	VoiceURL             string `url:"VoiceUrl,omitempty"`
}

// IncomingPhoneNumberCapabilities describes the capabilities afforded to a given IncomingPhoneNumber according to Twilio.
type IncomingPhoneNumberCapabilities struct {
	MMS   bool `json:"mms"`
	SMS   bool `json:"sms"`
	Voice bool `json:"voice"`
}

// IncomingPhoneNumber represents a phone number purchased from Twilio.
type IncomingPhoneNumber struct {
	AccountSID           string                          `json:"account_sid"`
	AddressRequirements  string                          `json:"address_requirements"`
	AddressSID           *string                         `json:"address_sid"`
	APIVersion           string                          `json:"api_version"`
	Beta                 bool                            `json:"beta"`
	BundleSID            *string                         `json:"bundle_sid"`
	Capabilities         IncomingPhoneNumberCapabilities `json:"capabilities"`
	DateCreated          dates.Rfc2822Time               `json:"date_created"`
	DateUpdated          dates.Rfc2822Time               `json:"date_updated"`
	EmergencyAddressSID  *string                         `json:"emergency_address_sid"`
	EmergencyStatus      string                          `json:"emergency_status"`
	FriendlyName         string                          `json:"friendly_name"`
	IdentitySID          *string                         `json:"identity_sid"`
	Origin               string                          `json:"origin"`
	PhoneNumber          string                          `json:"phone_number"`
	SID                  string                          `json:"sid"`
	SmsApplicationSID    *string                         `json:"sms_application_sid"`
	SmsFallbackMethod    string                          `json:"sms_fallback_method"`
	SmsFallbackURL       *string                         `json:"sms_fallback_url"`
	SmsMethod            string                          `json:"sms_method"`
	SmsURL               *string                         `json:"sms_url"`
	Status               string                          `json:"status"`
	StatusCallback       *string                         `json:"status_callback"`
	StatusCallbackMethod string                          `json:"status_callback_method"`
	TrunkSID             *string                         `json:"trunk_sid"`
	URI                  string                          `json:"uri"`
	VoiceApplicationSID  *string                         `json:"voice_application_sid"`
	VoiceCallerIDLookup  bool                            `json:"voice_caller_id_lookup"`
	VoiceFallbackMethod  string                          `json:"voice_fallback_method"`
	VoiceFallbackURL     *string                         `json:"voice_fallback_url"`
	VoiceMethod          string                          `json:"voice_method"`
	VoiceReceiveMode     string                          `json:"voice_receive_mode"`
	VoiceURL             *string                         `json:"voice_url"`
}

// CreateNewIncomingPhoneNumber purchases a new phone number in Twilio.
func (twilio *Twilio) CreateNewIncomingPhoneNumber(options CreateNewIncomingPhoneNumberOptions) (*IncomingPhoneNumber, error) {
	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	areaCode := params.Get("AreaCode")
	phoneNumber := params.Get("PhoneNumber")

	if areaCode == "" && phoneNumber == "" {
		return nil, errors.New("Missing required parameter PhoneNumber or AreaCode")
	}

	req, err := http.NewRequest(http.MethodPost, twilio.url("IncomingPhoneNumbers.json"), strings.NewReader(params.Encode()))

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

	response := new(IncomingPhoneNumber)

	decoder.Decode(&response)

	return response, nil
}

// DeleteIncomingPhoneNumber will release an existing IncomingPhoneNumber from Twilio.
func (twilio *Twilio) DeleteIncomingPhoneNumber(sid string) error {
	req, err := http.NewRequest(http.MethodDelete, twilio.url("IncomingPhoneNumbers/"+sid+".json"), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(twilio.credentials())

	res, err := twilio.delete(req)

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
