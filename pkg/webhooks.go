package twiligo

// SmsWebhook represents the response structure sent from Twilio for incoming SMS webhooks.
type SmsWebhook struct {
	AccountSID    string `json:"AccountSID"`
	APIVersion    string `json:"ApiVersion"`
	Body          string `json:"Body"`
	From          string `json:"From"`
	FromCity      string `json:"FromCity"`
	FromCountry   string `json:"FromCountry"`
	FromState     string `json:"FromState"`
	FromZip       string `json:"FromZip"`
	MessageSID    string `json:"MessageSid"`
	NumMedia      string `json:"NumMedia"`
	NumSegments   string `json:"NumSegments"`
	SmsMessageSID string `json:"SmsMessageSid"`
	SmsSID        string `json:"SmsSid"`
	SmsStatus     string `json:"SmsStatus"`
	To            string `json:"To"`
	ToCountry     string `json:"ToCountry"`
	ToCity        string `json:"ToCity"`
	ToState       string `json:"ToState"`
	ToZip         string `json:"ToZip"`
}
