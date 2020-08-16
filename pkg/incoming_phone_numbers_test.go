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

const createdIncomingPhoneNumberResponse = `{
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"address_requirements": "none",
	"address_sid": "ADXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"api_version": "2010-04-01",
	"beta": false,
	"capabilities": {
		"mms": true,
		"sms": true,
		"voice": true
	},
	"date_created": "Thu, 30 Jul 2020 00:00:00 +0000",
	"date_updated": "Thu, 30 Jul 2020 00:00:00 +0000",
	"emergency_status": "Active",
	"emergency_address_sid": "ADXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"friendly_name": "(555) 555-5555",
	"identity_sid": "RIXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"origin": "origin",
	"phone_number": "+15555555555",
	"sid": "PNXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"sms_application_sid": "",
	"sms_fallback_method": "POST",
	"sms_fallback_url": "",
	"sms_method": "POST",
	"sms_url": "",
	"status_callback": "",
	"status_callback_method": "POST",
	"trunk_sid": null,
	"uri": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/IncomingPhoneNumbers/PNXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX.json",
	"voice_application_sid": "",
	"voice_caller_id_lookup": false,
	"voice_fallback_method": "POST",
	"voice_fallback_url": null,
	"voice_method": "POST",
	"voice_url": null,
	"bundle_sid": "BUXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"voice_receive_mode": "voice",
	"status": "in-use"
}`

const errorCreatingIncomingPhoneNumberResponse = `{
	"code": 400,
	"message": "did or area_code is required",
	"more_info": "https://www.twilio.com/docs/errors/400",
	"status": 400
}`

const errorDeletingResourceResponse = `{
	"code": 20404,
	"message": "The request resource was not found",
	"more_info": "https://www.twilio.com/docs/errors/20404",
	"status": 404
}`

func TestWillMakeRequestToCreateNewIncomingNumberSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "IncomingPhoneNumbers.json"

		if strings.Contains(req.URL.Path, expected) == false {
			t.Logf("Incorrect URL supplied, expecting URL to contain [%s], but received [%s]", expected, req.URL.Path)
			t.Fail()
		}

		if req.Header.Get("Authorization") == "" {
			t.Log("Missing authorization credentials, they should be supplied via the Authorization header")
			t.Fail()
		}

		if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Logf("Incorrect content-type header supplied, expecting [%s], but received [%s]", "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
			t.Fail()
		}

		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("PhoneNumber") != "+15555555555" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Unique Name", params.Get("UniqueName"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdIncomingPhoneNumberResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewIncomingPhoneNumber(twiligo.CreateNewIncomingPhoneNumberOptions{
		PhoneNumber: "+15555555555",
	})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if response == nil {
		t.Log("Did not receive the expected response")
		t.Fail()
	}
}

func TestWillIncludeProperRequestBodyParametersWhenMakingRequestToCreateNewIncomingNumberSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("FriendlyName") != "(555) 555-5555" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "(555) 555-5555", params.Get("FriendlyName"))
			t.Fail()
		}

		if params.Get("PhoneNumber") != "+15555555555" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "+15555555555", params.Get("PhoneNumber"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdIncomingPhoneNumberResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	twilio.CreateNewIncomingPhoneNumber(twiligo.CreateNewIncomingPhoneNumberOptions{
		FriendlyName: "(555) 555-5555",
		PhoneNumber:  "+15555555555",
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToCreateNewIncomingNumber(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(errorCreatingIncomingPhoneNumberResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewIncomingPhoneNumber(twiligo.CreateNewIncomingPhoneNumberOptions{
		PhoneNumber: "+15555555555",
	})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "did or area_code is required"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}

func TestWillNotMakeRequestIfNotSuppliedWithOneOfTheRequiredParametersToCreateNewIncomingNumber(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		t.Logf("Request was incorrectly made, was not expecting the following request: %v", req)
		t.FailNow()

		return &http.Response{}
	})

	response, err := twilio.CreateNewIncomingPhoneNumber(twiligo.CreateNewIncomingPhoneNumberOptions{})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "Missing required parameter PhoneNumber or AreaCode"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}

func TestCanDeleteExistingIncomingPhoneNumberSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "IncomingPhoneNumbers/PN123456.json"

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

	err := twilio.DeleteIncomingPhoneNumber("PN123456")

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}
}

func TestWillHandleErrorResponsesWhenMakingRequestToDeleteExistingIncomingNumber(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(errorDeletingResourceResponse)),
			StatusCode: http.StatusNotFound,
			Header:     make(http.Header),
		}
	})

	err := twilio.DeleteIncomingPhoneNumber("PN123456")

	expected := "The request resource was not found"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
