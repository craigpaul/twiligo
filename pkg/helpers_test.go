package twiligo_test

import (
	"net/http"

	twiligo "github.com/craigpaul/twiligo/pkg"
)

type RoundTripFunc func(req *http.Request) *http.Response

func NewTestTwilio(fn RoundTripFunc) *twiligo.Twilio {
	client := &http.Client{
		Transport: RoundTripFunc(fn),
	}

	accountSID := "123"
	authToken := "456"

	twilio := twiligo.NewCustomClient(accountSID, authToken, client)

	return twilio
}

func (fn RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req), nil
}
