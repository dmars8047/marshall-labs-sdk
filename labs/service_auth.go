package labs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/dmars8047/marshall-labs-sdk/idam"
	"github.com/gin-gonic/gin"
)

const (
	ServiceAuthHeader = "X-Service-Auth"
	// serviceAuthMaxAge is the maximum age of a service auth header before it is rejected.
	serviceAuthMaxAge = 30 * time.Second
)

// serviceAuthPayload is the signed payload carried in the X-Service-Auth header.
type serviceAuthPayload struct {
	// Service is the name of the calling service.
	Service string `json:"service"`
	// IssuedAt is the UTC time the header was created.
	IssuedAt time.Time `json:"iat"`
}

// CreateServiceAuthHeader builds and signs an X-Service-Auth header value.
// Format: base64(JSON).base64(HMAC-SHA256(JSON))
func CreateServiceAuthHeader(serviceName string, key []byte) (string, error) {
	payload := serviceAuthPayload{
		Service:  serviceName,
		IssuedAt: time.Now().UTC(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(payloadBytes)
	sig := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(payloadBytes) + "." + base64.StdEncoding.EncodeToString(sig), nil
}

// VerifyServiceAuthHeader verifies the signature and freshness of an X-Service-Auth header value.
// Returns an error if the header is malformed, the signature is invalid, or the timestamp is stale.
func VerifyServiceAuthHeader(headerValue string, key []byte) error {
	parts := strings.SplitN(headerValue, ".", 2)
	if len(parts) != 2 {
		return errors.New("malformed service auth header")
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return errors.New("malformed service auth header: invalid payload encoding")
	}

	sigBytes, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return errors.New("malformed service auth header: invalid signature encoding")
	}

	mac := hmac.New(sha256.New, key)
	mac.Write(payloadBytes)
	expected := mac.Sum(nil)

	if !hmac.Equal(sigBytes, expected) {
		return errors.New("service auth header signature mismatch")
	}

	var payload serviceAuthPayload
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return errors.New("malformed service auth header: invalid payload")
	}

	age := time.Since(payload.IssuedAt)
	if age < 0 || age > serviceAuthMaxAge {
		return errors.New("service auth header is expired or from the future")
	}

	return nil
}

// ServiceAuthGuard returns a gin middleware that verifies the X-Service-Auth header.
// Requests missing or presenting an invalid header are rejected with 401 Unauthorized.
func ServiceAuthGuard(key []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		headerValue := c.GetHeader(ServiceAuthHeader)

		if headerValue == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, idam.NewErrorResponse(idam.InvalidRequestHeaders, idam.InvalidRequestHeadersMessage))
			return
		}

		if err := VerifyServiceAuthHeader(headerValue, key); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, idam.NewErrorResponse(idam.InvalidRequestHeaders, idam.InvalidRequestHeadersMessage))
			return
		}

		c.Next()
	}
}
