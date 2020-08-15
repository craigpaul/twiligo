package twiligo

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"sort"
)

// CheckSignature checks that the X-Twilio-Signature header on a request matches the expected signature defined by GenerateSignature.
func (twilio *Twilio) CheckSignature(r *http.Request, baseURL string) (bool, error) {
	err := r.ParseForm()

	if err != nil {
		return false, err
	}

	uri := baseURL + r.URL.String()

	expected, err := twilio.GenerateSignature(uri, r.PostForm)

	if err != nil {
		return false, err
	}

	actual := r.Header.Get("X-Twilio-Signature")

	if actual == "" {
		return false, errors.New("Request is missing an X-Twilio-Signature header")
	}

	return hmac.Equal(expected, []byte(actual)), nil
}

// GenerateSignature computes the Twilio signature for verifying the authenticity of a request. It is based on the specification at https://www.twilio.com/docs/security#validating-requests.
func (twilio *Twilio) GenerateSignature(uri string, values url.Values) ([]byte, error) {
	var buffer bytes.Buffer
	var expected bytes.Buffer

	buffer.WriteString(uri)

	keys := make(sort.StringSlice, 0, len(values))

	for key := range values {
		keys = append(keys, key)
	}

	keys.Sort()

	for _, key := range keys {
		buffer.WriteString(key)

		for _, value := range values[key] {
			buffer.WriteString(value)
		}
	}

	mac := hmac.New(sha1.New, []byte(twilio.AuthToken))
	mac.Write(buffer.Bytes())

	coder := base64.NewEncoder(base64.StdEncoding, &expected)

	_, err := coder.Write(mac.Sum(nil))

	if err != nil {
		return nil, err
	}

	err = coder.Close()

	if err != nil {
		return nil, err
	}

	return expected.Bytes(), nil
}
