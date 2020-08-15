package twiligo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	twiligo "github.com/craigpaul/twiligo/pkg"
)

const validSignatureParameters = `FromZip=89449&From=%2B15306666666&` +
	`FromCity=SOUTH+LAKE+TAHOE&ApiVersion=2010-04-01&To=%2B15306384866&` +
	`CallStatus=ringing&CalledState=CA&FromState=CA&Direction=inbound&` +
	`ToCity=OAKLAND&ToZip=94612&CallerCity=SOUTH+LAKE+TAHOE&FromCountry=US&` +
	`CallerName=CA+Wireless+Call&CalledCity=OAKLAND&CalledCountry=US&` +
	`Caller=%2B15306666666&CallerZip=89449&AccountSid=AC9a9f9392lad99kla0sklakjs90j092j3&` +
	`Called=%2B15306384866&CallerCountry=US&CalledZip=94612&CallSid=CAd800bb12c0426a7ea4230e492fef2a4f&` +
	`CallerState=CA&ToCountry=US&ToState=CA`

func TestCanCheckSignatureForGetRequest(t *testing.T) {
	authToken := "1c892n40nd03kdnc0112slzkl3091j20"
	twilio := twiligo.Twilio{AuthToken: authToken}

	uri, err := url.Parse("/1ed898x")

	if err != nil {
		t.Logf("Unexpected error occurred: %s", err)
		t.Fail()
	}

	headers := http.Header{
		"Content-Type":       []string{"application/x-www-form-urlencoded"},
		"X-Twilio-Signature": []string{"0JNGUhXSEKtzgrpd2HoQz7+F34M="},
	}

	req := http.Request{
		Method: "GET",
		URL:    uri,
		Header: headers,
	}

	query, err := url.ParseQuery(string(validSignatureParameters))

	if err != nil {
		t.Logf("Unexpected error occurred parsing the query parameters: %s", err)
		t.Fail()
	}

	req.URL.RawQuery = query.Encode()

	valid, err := twilio.CheckSignature(&req, "http://www.postbin.org")

	if err != nil {
		t.Logf("Unexpected error occurred while checking the signature: %s", err)
		t.Fail()
	}

	if !valid {
		t.Log("Expected signature to be valid, but was determined to be invalid")
		t.Fail()
	}

	headers["X-Twilio-Signature"] = []string{"foo"}

	valid, err = twilio.CheckSignature(&req, "http://www.postbin.org")

	if err != nil {
		t.Logf("Unexpected error occurred while checking the signature: %s", err)
		t.Fail()
	}

	if valid {
		t.Log("Expected signature to be invalid, but was determined to be valid")
		t.Fail()
	}

	delete(headers, "X-Twilio-Signature")

	_, err = twilio.CheckSignature(&req, "http://www.postbin.org")

	if err == nil {
		t.Log("Expected an error verifying that the given request is missing an X-Twilio-Signature header")
		t.Fail()
	}
}

func TestCanCheckSignatureForPostRequest(t *testing.T) {
	authToken := "1c892n40nd03kdnc0112slzkl3091j20"
	twilio := twiligo.Twilio{AuthToken: authToken}

	uri, err := url.Parse("/1ed898x")

	if err != nil {
		t.Logf("Unexpected error occurred: %s", err)
		t.Fail()
	}

	headers := http.Header{
		"Content-Type":       []string{"application/x-www-form-urlencoded"},
		"X-Twilio-Signature": []string{"fF+xx6dTinOaCdZ0aIeNkHr/ZAA="},
	}

	req := http.Request{
		Method: "POST",
		URL:    uri,
		Header: headers,
		Body:   ioutil.NopCloser(bytes.NewBufferString(validSignatureParameters)),
	}

	valid, err := twilio.CheckSignature(&req, "http://www.postbin.org")

	if err != nil {
		t.Logf("Unexpected error occurred while checking the signature: %s", err)
		t.Fail()
	}

	if !valid {
		t.Log("Expected signature to be valid, but was determined to be invalid")
		t.Fail()
	}

	headers["X-Twilio-Signature"] = []string{"foo"}

	valid, err = twilio.CheckSignature(&req, "http://www.postbin.org")

	if err != nil {
		t.Logf("Unexpected error occurred while checking the signature: %s", err)
		t.Fail()
	}

	if valid {
		t.Log("Expected signature to be invalid, but was determined to be valid")
		t.Fail()
	}

	delete(headers, "X-Twilio-Signature")

	_, err = twilio.CheckSignature(&req, "http://www.postbin.org")

	if err == nil {
		t.Log("Expected an error verifying that the given request is missing an X-Twilio-Signature header")
		t.Fail()
	}
}
