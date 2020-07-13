package twiligo

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

// This constant is used to represent a type of phone number to look up from Twilio.
const (
	Local PhoneNumberType = iota
	TollFree
	Mobile
)

// AvailablePhoneNumber represents a Twilio phone number that is currently available to be purchased.
type AvailablePhoneNumber struct {
	AddressRequirements string `json:"address_requirements"`
	Beta                bool   `json:"beta"`
	Capabilities        struct {
		MMS   bool `json:"mms"`
		SMS   bool `json:"sms"`
		Voice bool `json:"voice"`
	} `json:"capabilities"`
	FriendlyName string `json:"friendly_name"`
	IsoCountry   string `json:"iso_country"`
	LATA         string `json:"lata"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	PhoneNumber  string `json:"phone_number"`
	PostalCode   string `json:"postal_code"`
	RateCenter   string `json:"rate_center"`
	Region       string `json:"region"`
}

// AvailablePhoneNumbersResponse is the representation of the JSON response from Twilio when looking up available phone numbers.
type AvailablePhoneNumbersResponse struct {
	AvailablePhoneNumbers []*AvailablePhoneNumber `json:"available_phone_numbers"`
}

// GetAvailablePhoneNumberOptions are all of the options that can be provided to a GetAvailablePhoneNumbers call.
type GetAvailablePhoneNumberOptions struct {
	Page     int `url:",omitempty"`
	PageSize int `url:",omitempty"`
}

// PhoneNumberType is used to define whether a phone number is local, toll-free or mobile.
type PhoneNumberType int

// GetAvailablePhoneNumbers retrieves a listing of available phone numbers from Twilio.
func (twilio *Twilio) GetAvailablePhoneNumbers(country string, number PhoneNumberType, options GetAvailablePhoneNumberOptions) (*AvailablePhoneNumbersResponse, error) {
	resource := country + "/" + number.String()

	req, err := http.NewRequest(http.MethodGet, twilio.url("AvailablePhoneNumbers/"+resource+".json"), nil)

	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(twilio.credentials())

	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Encode()

	res, err := twilio.get(req)

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

	response := new(AvailablePhoneNumbersResponse)

	decoder.Decode(&response)

	return response, nil
}

func (number PhoneNumberType) String() string {
	return map[PhoneNumberType]string{
		Local:    "Local",
		TollFree: "TollFree",
		Mobile:   "Mobile",
	}[number]
}
