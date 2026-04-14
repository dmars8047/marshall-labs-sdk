package labs

import "time"

type VerifyTokenResponse struct {
	// The user ID parsed from the token
	UserId string `json:"user_id"`
	// The app ID parsed from the token
	AppId string `json:"app_id"`
	// The unique idetifier for the JWT token
	TokenId string `json:"token_id"`
	// UTC timestamp indicating when the auth token expires
	TokenExpiration time.Time `json:"token_expiration"`
}

func (response *VerifyTokenResponse) AsNewIdentityContext() *UserAuthContext {
	return &UserAuthContext{
		UserId:          response.UserId,
		ApplicationId:   response.AppId,
		TokenId:         response.TokenId,
		TokenExpiration: response.TokenExpiration,
	}
}
