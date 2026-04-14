package idam

import "github.com/dmars8047/strval"

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token         string `json:"token"`
	TokenType     string `json:"token_type"`
	ApplicationId string `json:"application"`
	ExpiresIn     int64  `json:"expires_in"`
	UserId        string `json:"user_id"`
	Username      string `json:"username"`
	RefreshToken  string `json:"refresh_token"`
}

func (request *UserLoginRequest) Validate() (valid bool, errors []string) {
	// Make suer the password and username wer passed in
	var validationErrors []string

	pwdValResult := strval.ValidateStringWithName(request.Password, "password", strval.MustNotBeEmpty())

	if !pwdValResult.Valid {
		validationErrors = append(validationErrors, pwdValResult.Messages...)
	}

	emailValResult := strval.ValidateStringWithName(request.Email, "email",
		strval.MustNotBeEmpty(), strval.MustBeValidEmailFormat())

	if !emailValResult.Valid {
		validationErrors = append(validationErrors, emailValResult.Messages...)
	}

	if len(validationErrors) > 0 {
		return false, validationErrors
	}

	return true, nil
}
