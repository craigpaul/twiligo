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

const createdSMSMessageResponse = `{
	"sid": "SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"date_created": "Thu, 30 Jul 2020 00:00:00 +0000",
	"date_updated": "Thu, 30 Jul 2020 00:00:00 +0000",
	"date_sent": null,
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"to": "+15555555555",
	"from": "+15555555554",
	"messaging_service_sid": null,
	"body": "Test Message",
	"status": "queued",
	"num_segments": "1",
	"num_media": "0",
	"direction": "outbound-api",
	"api_version": "2010-04-01",
	"price": null,
	"price_unit": "USD",
	"error_code": null,
	"error_message": null,
	"uri": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Messages/SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX.json",
	"subresource_uris": {
		"media": "/2010-04-01/Accounts/ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Messages/SMXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Media.json"
	}
}`

const errorCreatingNewSMSResponse = `{
	"code": 21602,
	"message": "Message body is required.",
	"more_info": "https://www.twilio.com/docs/errors/21602",
	"status": 400
}`

func TestWillMakeRequestToCreateNewSMSMessageSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "Messages.json"

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

		if params.Get("To") != "+15555555555" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "+15555555555", params.Get("To"))
			t.Fail()
		}

		if params.Get("Body") != "Test Message" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Test Message", params.Get("Body"))
			t.Fail()
		}

		if params.Get("From") != "+15555555554" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "+15555555554", params.Get("From"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdSMSMessageResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewSMSMessage("+15555555555", "Test Message", twiligo.CreateNewSMSMessageOptions{
		From: "+15555555554",
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

func TestWillIncludeProperRequestBodyParametersWhenMakingRequestToCreateNewSMSMessageSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("StatusCallback") != "https://example.com/callback" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "https://example.com/callback", params.Get("StatusCallback"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdSMSMessageResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	twilio.CreateNewSMSMessage("+15555555555", "Test Message", twiligo.CreateNewSMSMessageOptions{
		From:           "+15555555554",
		StatusCallback: "https://example.com/callback",
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToCreateNewSmsMessage(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(errorCreatingNewSMSResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewSMSMessage("+15555555555", "", twiligo.CreateNewSMSMessageOptions{})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "Message body is required."

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
