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

const createdChatUserResponse = `{
	"sid": "USXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"service_sid": "ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"role_sid": "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"identity": "Unique Name",
	"attributes": "{}",
	"is_online": null,
	"is_notifiable": null,
	"friendly_name": null,
	"joined_channels_count": 0,
	"date_created": "2020-07-30T00:00:00Z",
	"date_updated": "2020-07-30T00:00:00Z",
	"links": {
		"user_channels": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Users/USXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Channels",
		"user_bindings": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Users/USXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Bindings"
	},
	"url": "https://chat.twilio.com/v2/Services/ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Users/USXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
}`

const missingIdentityParameterResponse = `{
	"code": 20001,
	"message": "Missing required parameter Identity in the post body",
	"more_info": "https://www.twilio.com/docs/errors/20001",
	"status": 400
}`

func TestWillMakeRequestToCreateNewChatUserSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		expected := "Services/IS123/Users"

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

		if params.Get("Identity") != "Unique Name" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Unique Name", params.Get("Identity"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdChatUserResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewChatUser("Unique Name", "IS123", twiligo.CreateNewChatUserOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if response == nil {
		t.Log("Did not receive the expected response")
		t.Fail()
	}
}

func TestWillIncludeProperRequestBodyParametersWhenMakingRequestToCreateNewChatUserSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("RoleSid") != "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", params.Get("RoleSid"))
			t.Fail()
		}

		if params.Get("Attributes") != `{"custom":"value"}` {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", `{"custom":"value"}`, params.Get("Attributes"))
			t.Fail()
		}

		if params.Get("FriendlyName") != "Friendly Name" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Friendly Name", params.Get("FriendlyName"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdChatUserResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	twilio.CreateNewChatUser("Unique Name", "IS123", twiligo.CreateNewChatUserOptions{
		RoleSID:      "RLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		Attributes:   `{"custom":"value"}`,
		FriendlyName: "Friendly Name",
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToCreateNewChatUser(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(missingIdentityParameterResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewChatUser("", "IS123", twiligo.CreateNewChatUserOptions{})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "Missing required parameter Identity in the post body"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
