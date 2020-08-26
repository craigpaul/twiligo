package twiligo_test

import (
	"math/rand"
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

func RandomString(length int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	random := make([]rune, length)

	for index := range random {
		random[index] = letters[rand.Intn(len(letters))]
	}

	return string(random)
}
