package twiligo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	twiligo "github.com/craigpaul/twiligo/pkg"
)

const addedPhoneNumberToServiceResponse = `{
	"sid": "PNXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"service_sid": "KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"date_created": "2020-07-30T00:00:00Z",
	"date_updated": "2020-07-30T00:00:00Z",
	"phone_number": "+15555555555",
	"friendly_name": "(555) 555-5555",
	"iso_country": "CA",
	"capabilities": {
		"mms_inbound": true,
		"mms_outbound": true,
		"sms_inbound": true,
		"sms_outbound": true,
		"voice_inbound": true,
		"voice_outbound": true,
	},
	"url": "https://proxy.twilio.com/v1/Services/KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/PhoneNumbers/PNXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"is_reserved": false,
	"in_use": 0
}`

const alreadyAddedToServiceResponse = `{
	"code": 80104,
	"message": "PhoneNumber has already been added to Service",
	"more_info": "https://www.twilio.com/docs/errors/80104",
	"status": 400
}`

const errorRemovingPhoneNumberFromProxyServiceResponse = `{
	"code": 20404,
	"message": "The request resource was not found",
	"more_info": "https://www.twilio.com/docs/errors/20404",
	"status": 404
}`

func TestWillMakeRequestToAddPhoneNumberToExistingProxyServiceSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "Services/KS123/PhoneNumbers"

		if strings.Contains(req.URL.Path, expected) == false {
			t.Logf("Incorrect URL supplied, expecting URL to contain [%s], but received [%s]", expected, req.URL.Path)
			t.Fail()
		}

		if req.Header.Get("Authorization") == "" {
			t.Log("Missing authorization credentials, they should be supplied via the Authorization header")
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(addedPhoneNumberToServiceResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.AddPhoneNumberToProxyService("KS123", twiligo.AddPhoneNumberToProxyServiceOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if response == nil {
		t.Log("Did not receive the expected response")
		t.Fail()
	}
}

func TestWillIncludeProperRequestBodyParametersWhenMakingRequestToAddPhoneNumberToExistingProxyServiceSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("Sid") != "PN123" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "PN123", params.Get("Sid"))
			t.Fail()
		}

		if params.Get("PhoneNumber") != "+15555555555" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "+15555555555", params.Get("PhoneNumber"))
			t.Fail()
		}

		if params.Get("IsReserved") != "false" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "false", params.Get("IsReserved"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(addedPhoneNumberToServiceResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	reserved := false

	twilio.AddPhoneNumberToProxyService("KS123", twiligo.AddPhoneNumberToProxyServiceOptions{
		SID:         "PN123",
		PhoneNumber: "+15555555555",
		IsReserved:  &reserved,
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToAddPhoneNumberToExistingProxyServiceSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(alreadyAddedToServiceResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.AddPhoneNumberToProxyService("KS123", twiligo.AddPhoneNumberToProxyServiceOptions{
		SID: "PN123",
	})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "PhoneNumber has already been added to Service"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}

func TestCanRemovePhoneNumberFromProxyServiceSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "Services/KS123/PhoneNumbers/PN123"

		if strings.Contains(req.URL.Path, expected) == false {
			t.Logf("Incorrect URL supplied, expecting URL to contain [%s], but received [%s]", expected, req.URL.Path)
			t.Fail()
		}

		if req.Header.Get("Authorization") == "" {
			t.Log("Missing authorization credentials, they should be supplied via the Authorization header")
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(``)),
			StatusCode: http.StatusNoContent,
			Header:     make(http.Header),
		}
	})

	err := twilio.RemovePhoneNumberFromProxyService("KS123", "PN123")

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}
}

func TestWillHandleErrorResponsesWhenMakingRequestToRemovePhoneNumberFromProxyService(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(errorRemovingPhoneNumberFromProxyServiceResponse)),
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
		}
	})

	err := twilio.RemovePhoneNumberFromProxyService("KS123", "PN123")

	expected := "The request resource was not found"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
