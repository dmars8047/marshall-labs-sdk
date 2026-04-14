package labs

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dmars8047/marshall-labs-sdk/idam"
	"github.com/gin-gonic/gin"
)

const (
	UserAuthContextHeader = "X-Idam-User-Auth"
	authContextKey        = "idam-auth-context"
)

// Describes the identity and authentication context of a standard user.
type UserAuthContext struct {
	// ID of the application whos context was used during authentication.
	ApplicationId string `json:"application_id"`
	// ID of the user
	UserId string `json:"user_id"`
	// ID of the JWT Token
	TokenId string `json:"token_id"`
	// UTC timestamp indicating when the auth token expires
	TokenExpiration time.Time `json:"token_expiration"`
}

// Applies the user auth context properties to the required request headers.
// The header value is formatted as "<base64(json)>.<base64(hmac-sha256)>" so the
// recipient can verify integrity before trusting the payload.
func (c *UserAuthContext) ApplyRequestHeaders(request *http.Request, hmacKey []byte) error {
	contextBytes, err := json.Marshal(c)

	if err != nil {
		return err
	}

	mac := hmac.New(sha256.New, hmacKey)
	mac.Write(contextBytes)
	sig := mac.Sum(nil)

	headerValue := base64.StdEncoding.EncodeToString(contextBytes) + "." + base64.StdEncoding.EncodeToString(sig)
	request.Header.Add(UserAuthContextHeader, headerValue)

	return nil
}

// Middleware to add HeaderInfoProvider to Gin context.
// hmacKey must be the same key used by the caller when invoking ApplyRequestHeaders.
func UserAuthGuard(hmacKey []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("Parsing user auth context")

		authContextRaw := c.GetHeader(UserAuthContextHeader)

		if authContextRaw == "" {
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		parts := strings.SplitN(authContextRaw, ".", 2)

		if len(parts) != 2 {
			log.Printf("Malformed X-Idam-User-Auth header - %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		jsonBytes, err := base64.StdEncoding.DecodeString(parts[0])

		if err != nil {
			log.Printf("An error occurred when decoding idam user auth header - %s - %s", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		sigBytes, err := base64.StdEncoding.DecodeString(parts[1])

		if err != nil {
			log.Printf("An error occurred when decoding idam user auth signature - %s - %s", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		mac := hmac.New(sha256.New, hmacKey)
		mac.Write(jsonBytes)
		expectedSig := mac.Sum(nil)

		if !hmac.Equal(sigBytes, expectedSig) {
			log.Printf("Invalid X-Idam-User-Auth signature - %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		var context UserAuthContext

		err = json.Unmarshal(jsonBytes, &context)

		if err != nil {
			log.Printf("An error occurred when unmarshaling idam user auth header - %s - %s", c.Request.URL.Path, err)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		if context.TokenId == "" || context.UserId == "" || context.ApplicationId == "" {
			log.Printf("Invalid X-Idam-User-Auth header info - %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken, idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		if context.TokenExpiration.Before(time.Now().UTC()) {
			log.Printf("An expired token was used to access %s", c.Request.URL.Path)
			c.JSON(http.StatusUnauthorized,
				idam.NewErrorResponse(idam.InvalidAuthToken,
					idam.InvalidAuthTokenMessage))
			c.Abort()
			return
		}

		c.Set(authContextKey, &context)
		c.Next()
	}
}

// Get the AuthContext from the Gin context
func GetAuthContext(c *gin.Context) *UserAuthContext {
	return c.MustGet(authContextKey).(*UserAuthContext)
}
