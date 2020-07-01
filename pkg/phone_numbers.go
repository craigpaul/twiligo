package twiligo

import (
	"encoding/json"
	"net/http"

	"github.com/google/go-querystring/query"
)

// ...
const (
	Local PhoneNumberType = iota
	TollFree
	Mobile
)

// AvailablePhoneNumber ...
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

// AvailablePhoneNumbersResponse ...
type AvailablePhoneNumbersResponse struct {
	AvailablePhoneNumbers []*AvailablePhoneNumber `json:"available_phone_numbers"`
}

// GetAvailablePhoneNumberOptions ...
type GetAvailablePhoneNumberOptions struct {
	Page     int `url:",omitempty"`
	PageSize int `url:",omitempty"`
}

// PhoneNumberType ...
type PhoneNumberType int

// GetAvailablePhoneNumbers ...
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
