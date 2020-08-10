package twiligo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	twiligo "github.com/craigpaul/twiligo/pkg"
)

const okResponse = `{
	"available_phone_numbers": [
		{
			"address_requirements": "none",
			"beta": false,
			"capabilities": {
				"mms": false,
				"sms": true,
				"voice": true
			},
			"friendly_name": "(555) 555-5555",
			"iso_country": "CA",
			"lata": "888",
			"latitude":"52.133300",
			"locality": "Saskatoon",
			"longitude":"-106.666700",
			"phone_number": "+15555555555",
			"postal_code": "S7H 0S1",
			"rate_center": "SASKATOON",
			"region": "SK"
		}
	]
}`

const notFoundResponse = `{
	"status": 404,
	"message": "The requested resource was not found"	
}`

func TestWillMakeRequestToGetAvailableLocalPhoneNumbersSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "AvailablePhoneNumbers/CA/Local.json"

		if strings.Contains(req.URL.Path, expected) == false {
			t.Logf("Incorrect URL supplied, expecting URL to contain [%s], but received [%s]", expected, req.URL.Path)
			t.Fail()
		}

		if req.Header.Get("Authorization") == "" {
			t.Log("Missing authorization credentials, they should be supplied via the Authorization header")
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(okResponse)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		}
	})

	numbers, err := twilio.GetAvailablePhoneNumbers("CA", twiligo.Local, twiligo.GetAvailablePhoneNumberOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if len(numbers) == 0 {
		t.Log("Did not receive the expected available phone number in the response")
		t.Fail()
	}
}

func TestWillMakeRequestToGetAvailableTollFreePhoneNumbersSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "AvailablePhoneNumbers/CA/TollFree.json"

		if strings.Contains(req.URL.Path, expected) == false {
			t.Logf("Incorrect URL supplied, expecting URL to contain [%s], but received [%s]", expected, req.URL.Path)
			t.Fail()
		}

		if req.Header.Get("Authorization") == "" {
			t.Log("Missing authorization credentials, they should be supplied via the Authorization header")
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(okResponse)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		}
	})

	numbers, err := twilio.GetAvailablePhoneNumbers("CA", twiligo.TollFree, twiligo.GetAvailablePhoneNumberOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if len(numbers) == 0 {
		t.Log("Did not receive the expected available phone number in the response")
		t.Fail()
	}
}

func TestWillMakeRequestToGetAvailableMobilePhoneNumbersSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "AvailablePhoneNumbers/CA/Mobile.json"

		if strings.Contains(req.URL.Path, expected) == false {
			t.Logf("Incorrect URL supplied, expecting URL to contain [%s], but received [%s]", expected, req.URL.Path)
			t.Fail()
		}

		if req.Header.Get("Authorization") == "" {
			t.Log("Missing authorization credentials, they should be supplied via the Authorization header")
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(okResponse)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		}
	})

	numbers, err := twilio.GetAvailablePhoneNumbers("CA", twiligo.Mobile, twiligo.GetAvailablePhoneNumberOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if len(numbers) == 0 {
		t.Log("Did not receive the expected available phone number in the response")
		t.Fail()
	}
}

func TestWillIncludeProperQueryParametersWhenMakingRequestToGetAvailablePhoneNumbersSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		params := req.URL.Query()

		if params.Get("MmsEnabled") != "false" {
			t.Logf("Incorrect query parameter supplied, expecting [%s], but received [%s]", "false", params.Get("MmsEnabled"))
			t.Fail()
		}

		if params.Get("Page") != "2" {
			t.Logf("Incorrect query parameter supplied, expecting [%s], but received [%s]", "2", params.Get("Page"))
			t.Fail()
		}

		if params.Get("PageSize") != "50" {
			t.Logf("Incorrect query parameter supplied, expecting [%s], but received [%s]", "50", params.Get("PageSize"))
			t.Fail()
		}

		if params.Get("SmsEnabled") != "true" {
			t.Logf("Incorrect query parameter supplied, expecting [%s], but received [%s]", "true", params.Get("SmsEnabled"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(okResponse)),
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
		}
	})

	mmsEnabled := false
	smsEnabled := true

	twilio.GetAvailablePhoneNumbers("CA", twiligo.Mobile, twiligo.GetAvailablePhoneNumberOptions{
		MmsEnabled: &mmsEnabled,
		Page:       2,
		PageSize:   50,
		SmsEnabled: &smsEnabled,
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToGetAvailablePhoneNumbers(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(notFoundResponse)),
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
		}
	})

	numbers, err := twilio.GetAvailablePhoneNumbers("CA", twiligo.Mobile, twiligo.GetAvailablePhoneNumberOptions{})

	if numbers != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", numbers)
		t.Fail()
	}

	expected := "The requested resource was not found"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}

func TestWillConvertGivenNumberTypeAndCountryToMatchingPhoneNumberTypeToPrevent404ErrorsThroughTwilioAPI(t *testing.T) {
	cases := map[string]map[twiligo.PhoneNumberType]twiligo.PhoneNumberType{
		"CA": {twiligo.Mobile: twiligo.Local},
		"US": {twiligo.Mobile: twiligo.Local},
		"FR": {twiligo.Mobile: twiligo.Mobile, twiligo.Local: twiligo.Local},
	}

	for country, values := range cases {
		for given, expected := range values {
			numberType := twiligo.GetPhoneNumberType(int(given), country)

			if numberType != expected {
				t.Logf("Incorrect number type returned, expected [%s], but received [%s]", expected, numberType)
			}
		}
	}
}
