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

const createdResponse = `{
	"sid": "KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"account_sid": "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"chat_instance_sid": null,
	"unique_name": "Unique Name",
	"default_ttl": 0,
	"callback_url": null,
	"geo_match_level": "country",
	"number_selection_behavior": "prefer-sticky",
	"intercept_callback_url": null,
	"out_of_session_callback_url": null,
	"date_created": "2020-07-30T00:00:00Z",
	"date_updated": "2020-07-30T00:00:00Z",
	"url": "https://proxy.twilio.com/v1/Services/KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	"links": {
		"sessions": "https://proxy.twilio.com/v1/Services/KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/Sessions",
		"phone_numbers": "https://proxy.twilio.com/v1/Services/KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/PhoneNumbers",
		"short_codes": "https://proxy.twilio.com/v1/Services/KSXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX/ShortCodes"
	}
}`

const missingNameResponse = `{
	"code": 20001,
	"message": "Missing required parameter UniqueName in the post body",
	"more_info": "https://www.twilio.com/docs/errors/20001",
	"status": 400
}`

func TestWillMakeRequestToCreateNewProxyServiceSuccessfully(t *testing.T) {
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

		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("UniqueName") != "Unique Name" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "Unique Name", params.Get("UniqueName"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewProxyService("Unique Name", twiligo.CreateNewProxyServiceOptions{})

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	if response == nil {
		t.Log("Did not receive the expected response")
		t.Fail()
	}
}

func TestWillIncludeProperRequestBodyParametersWhenMakingRequestToCreateNewProxyServiceSuccessfully(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		body, _ := ioutil.ReadAll(req.Body)
		params, _ := url.ParseQuery(string(body))

		if params.Get("DefaultTtl") != "3600" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "3600", params.Get("DefaultTtl"))
			t.Fail()
		}

		if params.Get("GeoMatchLevel") != "area-code" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "area-code", params.Get("GeoMatchLevel"))
			t.Fail()
		}

		if params.Get("NumberSelectionBehaviour") != "avoid-sticky" {
			t.Logf("Incorrect request parameter supplied, expecting [%s], but received [%s]", "avoid-sticky", params.Get("NumberSelectionBehaviour"))
			t.Fail()
		}

		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(createdResponse)),
			StatusCode: http.StatusCreated,
			Header:     make(http.Header),
		}
	})

	twilio.CreateNewProxyService("Unique Name", twiligo.CreateNewProxyServiceOptions{
		DefaultTTL:               3600,
		GeoMatchLevel:            twiligo.AreaCode,
		NumberSelectionBehaviour: twiligo.AvoidSticky,
	})
}

func TestWillHandleErrorResponsesWhenMakingRequestToCreateNewProxyService(t *testing.T) {
	twilio := NewTestTwilio(func(req *http.Request) *http.Response {
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString(missingNameResponse)),
			StatusCode: http.StatusBadRequest,
			Header:     make(http.Header),
		}
	})

	response, err := twilio.CreateNewProxyService("", twiligo.CreateNewProxyServiceOptions{})

	if response != nil {
		t.Logf("Response was incorrectly returned, was not expecting the following response: %v", response)
		t.Fail()
	}

	expected := "Missing required parameter UniqueName in the post body"

	if err.Error() != expected {
		t.Logf("Incorrect error returned, expected [%s], but received [%s]", expected, err)
		t.Fail()
	}
}
