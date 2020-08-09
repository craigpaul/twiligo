package twiligo_test

import (
	"fmt"
	"strconv"
	"testing"

	twiligo "github.com/craigpaul/twiligo/pkg"
	"github.com/dgrijalva/jwt-go"
)

const accountSID = "ACXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
const signingKeySID = "SKXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
const ttl = 3600

func TestCanGenerateAccessTokenWithEmptyGrants(t *testing.T) {
	identity := "unique_identity"
	secret := RandomString(32)

	token := twiligo.NewAccessToken(twiligo.NewAccessTokenOptions{
		AccountSID:    accountSID,
		Identity:      &identity,
		SigningKeySID: signingKeySID,
		Secret:        secret,
		TTL:           ttl,
	})

	jwtToken, err := token.ToJWT()

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	tok, err := jwt.Parse(jwtToken, getKeyFunc(secret))

	if tok.Valid == false {
		t.Log("Expecting a valid token to be parsed, but it was determined to be invalid")
		t.Fail()
	}

	validateClaims(t, tok)
}

func TestCanGenerateAccessTokenWithChatGrant(t *testing.T) {
	identity := "unique_identity"
	secret := RandomString(32)

	token := twiligo.NewAccessToken(twiligo.NewAccessTokenOptions{
		AccountSID:    accountSID,
		Identity:      &identity,
		SigningKeySID: signingKeySID,
		Secret:        secret,
		TTL:           ttl,
	})

	serviceSID := "ISXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

	grant := twiligo.NewChatGrant(twiligo.NewChatGrantOptions{
		ServiceSID: &serviceSID,
	})

	token.AddGrant(grant)

	jwtToken, err := token.ToJWT()

	if err != nil {
		t.Logf("Error was incorrectly returned, was not expecting the following error: %s", err)
		t.Fail()
	}

	tok, err := jwt.Parse(jwtToken, getKeyFunc(secret))

	if tok.Valid == false {
		t.Log("Expecting a valid token to be parsed, but it was determined to be invalid")
		t.Fail()
	}

	validateClaims(t, tok)

	claims := tok.Claims.(jwt.MapClaims)
	grants := claims["grants"].(map[string]interface{})

	if grants["chat"] == nil {
		t.Log("Expecting a chat grant to exist, but it was not found")
		t.Fail()
	}
}

func getKeyFunc(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	}
}

func validateClaims(t *testing.T, token *jwt.Token) {
	claims := token.Claims.(jwt.MapClaims)

	sub := claims["sub"]
	exp := int64(claims["exp"].(float64))
	iat := int64(claims["iat"].(float64))
	iss := claims["iss"]
	jti := claims["jti"]
	grants := claims["grants"].(map[string]interface{})
	identity := grants["identity"]

	if sub != accountSID {
		t.Logf("Incorrect sub claim returned, expecting [%s], but received [%s]", accountSID, sub)
		t.Fail()
	}

	if iat+ttl != exp {
		t.Logf("Incorrect expiry claim returned, expecting [%d], but received [%d]", exp, iat+ttl)
		t.Fail()
	}

	if iss != signingKeySID {
		t.Logf("Incorrect issuer claim returned, expecting [%s], but received [%s]", signingKeySID, iss)
		t.Fail()
	}

	expected := signingKeySID + "-" + strconv.FormatInt(iat, 10)

	if jti.(string) != expected {
		t.Logf("Incorrect JWT ID supplied, expecting [%s], but received [%s]", expected, jti.(string))
		t.Fail()
	}

	if identity != "unique_identity" {
		t.Logf("Incorrect identity supplied, expecting [%s], but received [%s]", "unique_identity", identity)
		t.Fail()
	}
}
