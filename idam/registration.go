package idam

import (
	"time"

	"github.com/dmars8047/strval"
)

type ExistingUserRegistrationRequest struct {
	Email string `json:"email"`
}

type UserRegistrationRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	InviteToken string `json:"invite_token,omitempty"`
}

type UserRegistrationResponse struct {
	UserId       string    `json:"user_id"`
	Email        string    `json:"email"`
	Verified     bool      `json:"verified"`
	Provider     string    `json:"provider"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
	Features     []string  `json:"features"`
}

func (request *UserRegistrationRequest) Validate() (valid bool, errors []string) {
	validationErrors := make([]string, 0)

	emailValidationResult := strval.ValidateStringWithName(request.Email, "email",
		strval.MustNotBeEmpty(),
		strval.MustBeValidEmailFormat(),
		strval.MustHaveMaxLengthOf(254))

	if !emailValidationResult.Valid {
		validationErrors = append(validationErrors, emailValidationResult.Messages...)
	}

	passwordValidationResult := validatePassword(request.Password, "password")

	if !passwordValidationResult.Valid {
		validationErrors = append(validationErrors, passwordValidationResult.Messages...)
	}

	if len(validationErrors) > 0 {
		return false, validationErrors
	}

	return true, nil
}
