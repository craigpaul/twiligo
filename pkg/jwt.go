package twiligo

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AccessToken ...
type AccessToken struct {
	AccountSID    string
	Grants        []Grant
	Identity      *string
	SigningKeySID string
	Secret        string
	TTL           int
}

// ChatGrant ...
type ChatGrant struct {
	ServiceSID *string
}

// Grant ...
type Grant interface {
	GetKey() string
	GetPayload() map[string]string
}

// NewAccessTokenOptions ...
type NewAccessTokenOptions struct {
	AccountSID    string
	Identity      *string
	SigningKeySID string
	Secret        string
	TTL           int
}

// NewChatGrantOptions ...
type NewChatGrantOptions struct {
	ServiceSID *string
}

// NewAccessToken ...
func NewAccessToken(options NewAccessTokenOptions) AccessToken {
	return AccessToken{
		AccountSID:    options.AccountSID,
		Identity:      options.Identity,
		SigningKeySID: options.SigningKeySID,
		Secret:        options.Secret,
		TTL:           options.TTL,
	}
}

// NewChatGrant ...
func NewChatGrant(options NewChatGrantOptions) ChatGrant {
	return ChatGrant{
		ServiceSID: options.ServiceSID,
	}
}

// AddGrant ...
func (token *AccessToken) AddGrant(grant Grant) {
	token.Grants = append(token.Grants, grant)
}

// ToJWT ...
func (token *AccessToken) ToJWT() (string, error) {
	now := time.Now()
	method := jwt.SigningMethodHS256
	signingKey := []byte(token.Secret)

	grants := make(map[string]interface{})

	for _, grant := range token.Grants {
		payload := grant.GetPayload()

		if payload == nil {
			payload = make(map[string]string)
		}

		grants[grant.GetKey()] = payload
	}

	if token.Identity != nil {
		grants["identity"] = token.Identity
	}

	jsonWebToken := jwt.NewWithClaims(method, jwt.MapClaims{
		"jti":    token.SigningKeySID + "-" + strconv.FormatInt(now.Unix(), 10),
		"iss":    token.SigningKeySID,
		"sub":    token.AccountSID,
		"iat":    now.Unix(),
		"exp":    now.Add(time.Second * time.Duration(token.TTL)).Unix(),
		"grants": grants,
	})

	jsonWebToken.Header = map[string]interface{}{
		"alg": method.Name,
		"cty": "twilio-fpa;v=1",
		"typ": "JWT",
	}

	signedString, err := jsonWebToken.SignedString(signingKey)

	if err != nil {
		return "", err
	}

	return signedString, nil
}

func (token *AccessToken) String() string {
	t, err := token.ToJWT()

	if err != nil {
		return ""
	}

	return t
}

// GetKey ...
func (grant ChatGrant) GetKey() string {
	return "chat"
}

// GetPayload ...
func (grant ChatGrant) GetPayload() map[string]string {
	payload := make(map[string]string)

	if grant.ServiceSID != nil {
		payload["service_sid"] = *grant.ServiceSID
	}

	return payload
}
