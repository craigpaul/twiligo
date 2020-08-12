package twiligo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

const createdChatServiceResponse = `{
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"date_created": "2020-07-30T00:00:00Z",
	"date_updated": "2020-07-30T00:00:00Z",
	"default_channel_creator_role_sid": "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"default_channel_role_sid": "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"default_service_role_sid": "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"friendly_name": "Friendly Name",
	"limits": {
		"user_channels": 250,
		"channel_members": 100
	},
	"links": {
		"channels": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Channels",
		"bindings": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Bindings",
		"users": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Users",
		"roles": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Roles"
	},
	"notifications": {
		"removed_from_channel": {
			"enabled": false
		},
		"log_enabled": false,
		"added_to_channel": {
			"enabled": false
		},
		"new_message": {
			"enabled": false
		},
		"invited_to_channel": {
			"enabled": false
		}
	},
	"media": {
		"compatibility_message": "Media messages are not supported by your client",
		"size_limit_mb": 150
	},
	"post_webhook_url": null,
	"pre_webhook_url": null,
	"pre_webhook_retry_count": 0,
	"post_webhook_retry_count": 0,
	"reachability_enabled": false,
	"read_status_enabled": true,
	"sid": "ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"typing_indicator_timeout": 5,
	"url": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"webhook_filters": null,
	"webhook_method": null
}`

const missingChatServiceNameResponse = `{
	"code": 20001,
	"message": "Missing required parameter FriendlyName in the post body",
	"more_info": "https://www.twilio.com/docs/errors/20001",
	"status": 400
}`

func TestWillMakeRequestToCreateNewChatServiceSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "Services"

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

		if params.Get("FriendlyName") != "Friendly Name" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Friendly Name", params.Get("FriendlyName"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdChatServiceResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewChatService("Friendly Name")

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if response == nil {
		t.Log("Did not receive the expected response")
		t.Fail()
	}
}

func TestWillHandleErrorResponsesWhenMakingRequestToCreateNewChatService(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(missingChatServiceNameResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewChatService("")

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "Missing required parameter FriendlyName in the post body"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
