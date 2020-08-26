package twiligo_test

import (
	"testing"

	twiligo "github.com/craigpaul/twiligo/pkg"
)

func TestWillConvertGivenDirectionStringToMatchingMessageDirection(t *testing.T) {
	cases := map[string]twiligo.MessageDirection{
		"inbound":        twiligo.Inbound,
		"outbound-api":   twiligo.OutboundAPI,
		"outbound-call":  twiligo.OutboundCall,
		"outbound-reply": twiligo.OutboundReply,
	}

	for given, expected := range cases {
		direction := twiligo.ConvertDirectionToMessageDirection(given)

		if direction != expected {
			t.Logf("Incorrect message direction returned, expected [%s], but received [%s]", expected, direction)
			t.Fail()
		}
	}
}

func TestWillConvertGivenStatusStringToMatchingMessageStatus(t *testing.T) {
	cases := map[string]twiligo.MessageStatus{
		"accepted":    twiligo.Accepted,
		"queued":      twiligo.Queued,
		"sending":     twiligo.Sending,
		"sent":        twiligo.Sent,
		"failed":      twiligo.Failed,
		"delivered":   twiligo.Delivered,
		"undelivered": twiligo.Undelivered,
		"receiving":   twiligo.Receiving,
		"received":    twiligo.Received,
	}

	for given, expected := range cases {
		status := twiligo.ConvertStatusToMessageStatus(given)

		if status != expected {
			t.Logf("Incorrect message status returned, expected [%s], but received [%s]", expected, status)
			t.Fail()
		}
	}
}

func TestWillConvertGivenNumberTypeAndCountryToMatchingPhoneNumberTypeToPrevent404ErrorsThroughTwilioAPI(t *testing.T) {
	cases := map[string]map[twiligo.PhoneNumberType]twiligo.PhoneNumberType{
		"CA": {twiligo.Mobile: twiligo.Local},
		"US": {twiligo.Mobile: twiligo.Local},
		"FR": {twiligo.Mobile: twiligo.Mobile, twiligo.Local: twiligo.Local},
	}

	for country, values := range cases {
		for given, expected := range values {
			numberType := twiligo.GetPhoneNumberType(int(given), country)

			if numberType != expected {
				t.Logf("Incorrect number type returned, expected [%s], but received [%s]", expected, numberType)
				t.Fail()
			}
		}
	}
}
