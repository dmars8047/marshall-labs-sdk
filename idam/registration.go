package idam

import (
	"time"

	"github.com/dmars8047/strval"
)

type ExistingUserRegistrationRequest struct {
	Email string `json:"email"`
}

type UserRegistrationRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRegistrationResponse struct {
	UserId       string    `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Verified     bool      `json:"verified"`
	Provider     string    `json:"provider"`
	CreatedAtUTC time.Time `json:"created_at_utc"`
	Features     []string  `json:"features"`
}

// Registration returns a User object if the registration was successful

func (request *UserRegistrationRequest) Validate() (valid bool, errors []string) {
	validationErrors := make([]string, 0)

	// Validate the username
	// The username must be alphanumeric, be at least 3 characters long, and have a max length of 20 characters
	usrnameValidationResult := strval.ValidateStringWithName(request.Username, "username",
		strval.MustNotBeEmpty(),
		strval.MustBeAlphaNumeric(),
		strval.MustHaveMinLengthOf(3),
		strval.MustHaveMaxLengthOf(20))

	if !usrnameValidationResult.Valid {
		validationErrors = append(validationErrors, usrnameValidationResult.Messages...)
	}

	// Validate email
	// The email must be not empty and valid email address
	emailValidationResult := strval.ValidateStringWithName(request.Email, "email",
		strval.MustNotBeEmpty(),
		strval.MustBeValidEmailFormat(),
		strval.MustHaveMaxLengthOf(254))

	if !emailValidationResult.Valid {
		validationErrors = append(validationErrors, emailValidationResult.Messages...)
	}

	// The password must contain at least one special character, number, uppercase letter, and lowercase letter
	// The password must be at least 8 characters long and have a max length of 64 characters
	// The password must contain at least one of the following special characters: !@#$%^&*()_+-={}[]|\:,.?
	// The password must contain only printable characters and ascii characters
	// The password must not contain any of the following special characters: '"<>`~;
	passwordValidationResult := strval.ValidateStringWithName(request.Password, "password",
		strval.MustNotBeEmpty(),
		strval.MustHaveMinLengthOf(MinPasswordLength),
		strval.MustHaveMaxLengthOf(MaxPasswordLength),
		strval.MustContainAtLeastOne([]rune(AllowablePasswordSpecialCharacters)),
		strval.MustNotContainAnyOf([]rune(DisallowedPassowrdSpecialCharacters)),
		strval.MustContainNumbers(),
		strval.MustContainUppercaseLetter(),
		strval.MustContainLowercaseLetter(),
		strval.MustOnlyContainPrintableCharacters(),
		strval.MustOnlyContainASCIICharacters())

	if !passwordValidationResult.Valid {
		validationErrors = append(validationErrors, passwordValidationResult.Messages...)
	}

	if len(validationErrors) > 0 {
		return false, validationErrors
	}

	return true, nil
}
