package twiligo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	twiligo "github.com/craigpaul/twiligo/pkg"
)

const attributesAreJsonResponse = `{
	"code": 50304,
	"message": "Attributes not valid JSON",
	"more_info": "https://www.twilio.com/docs/errors/50304",
	"status": 400
}`

const createdConversationResponse = `{
	"date_updated": "2020-07-30T00:00:00Z",
	"friendly_name": "Friendly Conversation",
	"timers": {},
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"url": "https://conversations.twilio.com/v1/Conversations/CHXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"state": "active",
	"date_created": "2020-07-30T00:00:00Z",
	"messaging_service_sid": "MGXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"sid": "CHXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"attributes": "{}",
	"chat_service_sid": "ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"links": {
		"participants": "https://conversations.twilio.com/v1/Conversations/CHXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Participants",
		"messages": "https://conversations.twilio.com/v1/Conversations/CHXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Messages",
		"webhooks": "https://conversations.twilio.com/v1/Conversations/CHXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Webhooks"
	}
}`

func TestWillMakeRequestToCreateNewConversationSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "Conversations"

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

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdConversationResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewConversation(twiligo.CreateNewConversationOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if response == nil {
		t.Log("Did not receive the expected response")
		t.Fail()
	}
}

func TestWillIncludeProperRequestBodyParametersWhenMakingRequestToCreateNewConversationSuccessfully(t *testing.T) {
	now := time.Now()
	attributes := "{\"custom\":\"value\"}"

	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("FriendlyName") != "Friendly Conversation" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Friendly Conversation", params.Get("FriendlyName"))
			t.Fail()
		}

		if params.Get("DateCreated") != now.Format(time.RFC3339) {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", now.Format(time.RFC3339), params.Get("DateCreated"))
			t.Fail()
		}

		if params.Get("DateUpdated") != now.Format(time.RFC3339) {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", now.Format(time.RFC3339), params.Get("DateUpdated"))
			t.Fail()
		}

		if params.Get("Attributes") != attributes {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", attributes, params.Get("Attributes"))
			t.Fail()
		}

		if params.Get("State") != "active" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "active", params.Get("State"))
			t.Fail()
		}

		if params.Get("Timers.Inactive") != "PT10M" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "PT10M", params.Get("Timers.Inactive"))
			t.Fail()
		}

		if params.Get("Timers.Closed") != "PT1H30M" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "PT1H30M", params.Get("Timers.Closed"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdConversationResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	twilio.CreateNewConversation(twiligo.CreateNewConversationOptions{
		FriendlyName:  "Friendly Conversation",
		DateCreated:   now,
		DateUpdated:   now,
		Attributes:    attributes,
		State:         "active",
		InactiveTimer: "PT10M",
		ClosedTimer:   "PT1H30M",
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToCreateNewConversation(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(attributesAreJsonResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewConversation(twiligo.CreateNewConversationOptions{
		Attributes: "{\"custom\":\"value}",
	})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "Attributes not valid JSON"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
