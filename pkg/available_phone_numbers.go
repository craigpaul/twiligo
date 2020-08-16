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
	Locality     string `json:"locality"`
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
	AreaCode                      int    `url:",omitempty"`
	Beta                          *bool  `url:",omitempty"`
	Contains                      string `url:",omitempty"`
	Distance                      int    `url:",omitempty"`
	ExcludeAllAddressRequired     *bool  `url:",omitempty"`
	ExcludeForeignAddressRequired *bool  `url:",omitempty"`
	ExcludeLocalAddressRequired   *bool  `url:",omitempty"`
	InLata                        string `url:",omitempty"`
	InLocality                    string `url:",omitempty"`
	InPostalCode                  string `url:",omitempty"`
	InRateCenter                  string `url:",omitempty"`
	InRegion                      string `url:",omitempty"`
	MmsEnabled                    *bool  `url:",omitempty"`
	NearLatLong                   string `url:",omitempty"`
	NearNumber                    string `url:",omitempty"`
	Page                          int    `url:",omitempty"`
	PageSize                      int    `url:",omitempty"`
	SmsEnabled                    *bool  `url:",omitempty"`
	VoiceEnabled                  *bool  `url:",omitempty"`
}

// PhoneNumberType is used to define whether a phone number is local, toll-free or mobile.
type PhoneNumberType int

// GetAvailablePhoneNumbers retrieves a listing of available phone numbers from Twilio.
func (twilio *Twilio) GetAvailablePhoneNumbers(country string, number PhoneNumberType, options GetAvailablePhoneNumberOptions) ([]*AvailablePhoneNumber, error) {
	resource := country + "/" + number.String()

	params, err := query.Values(options)

	if err != nil {
		return nil, err
	}

	res, err := twilio.get(twilio.url("AvailablePhoneNumbers/"+resource+".json"), &params)

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

	return response.AvailablePhoneNumbers, nil
}

// GetPhoneNumberType will convert a given integer into a PhoneNumberType by the given country. Certain countries do not support all PhoneNumberType values, so this function can be used as a safe mapping based on the values from this document https://support.twilio.com/hc/en-us/articles/223183068-Twilio-international-phone-number-availability-and-their-capabilities. Note: Not all cases are currently supported, but can be amended as necessary.
func GetPhoneNumberType(number int, country string) PhoneNumberType {
	numberType := PhoneNumberType(number)

	if (country == "CA" || country == "US") && numberType == Mobile {
		numberType = Local
	}

	return numberType
}

func (number PhoneNumberType) String() string {
	return map[PhoneNumberType]string{
		Local:    "Local",
		TollFree: "TollFree",
		Mobile:   "Mobile",
	}[number]
}
