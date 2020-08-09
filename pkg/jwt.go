package twiligo

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// AccessToken holds the necessary important information for encoding a JWT to use with various Twilio services.
type AccessToken struct {
	AccountSID    string
	Grants        []Grant
	Identity      *string
	SigningKeySID string
	Secret        string
	TTL           int
}

// ChatGrant represents a specific type of Grant that is necessary to use Twilio's Programmable Chat.
type ChatGrant struct {
	ServiceSID *string
}

// Grant represents a type of permission that can be attached to an AccessToken.
type Grant interface {
	GetKey() string
	GetPayload() map[string]string
}

// NewAccessTokenOptions represents the possible options that can be provided to a new AccessToken.
type NewAccessTokenOptions struct {
	AccountSID    string
	Identity      *string
	SigningKeySID string
	Secret        string
	TTL           int
}

// NewChatGrantOptions represents the possible options that can be provided to a new ChatGrant.
type NewChatGrantOptions struct {
	ServiceSID *string
}

// NewAccessToken creates a new AccessToken with the given options.
func NewAccessToken(options NewAccessTokenOptions) AccessToken {
	return AccessToken{
		AccountSID:    options.AccountSID,
		Identity:      options.Identity,
		SigningKeySID: options.SigningKeySID,
		Secret:        options.Secret,
		TTL:           options.TTL,
	}
}

// NewChatGrant creates a new ChatGrant with the given options.
func NewChatGrant(options NewChatGrantOptions) ChatGrant {
	return ChatGrant{
		ServiceSID: options.ServiceSID,
	}
}

// AddGrant appends a given grant to the current AccessToken.
func (token *AccessToken) AddGrant(grant Grant) {
	token.Grants = append(token.Grants, grant)
}

// ToJWT converts the given AccessToken attributes into a properly encoded and signed JWT.
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

// GetKey returns the string identifier for the current grant.
func (grant ChatGrant) GetKey() string {
	return "chat"
}

// GetPayload generates the full custom payload for the current grant.
func (grant ChatGrant) GetPayload() map[string]string {
	payload := make(map[string]string)

	if grant.ServiceSID != nil {
		payload["service_sid"] = *grant.ServiceSID
	}

	return payload
}
