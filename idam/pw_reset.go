package idam

import "github.com/dmars8047/strval"

// UserPasswordResetExecutionRequest is the request object for the password reset execution endpoint
type UserPasswordResetExecutionRequest struct {
	// The user id of the user to reset the password for
	UserID string `json:"user_id"`
	// The new password for the user
	NewPassword string `json:"new_password"`
	// The reset token for the user
	PasswordResetToken string `json:"password_reset_token"`
	// The verification code for the user
	VerificationCode string `json:"verification_code"`
}

// Validate validates the password reset execution request
func (request *UserPasswordResetExecutionRequest) Validate() (valid bool, errors []string) {
	// The password must contain at least one special character, number, uppercase letter, and lowercase letter
	// The password must be at least 8 characters long and have a max length of 64 characters
	// The password must contain at least one of the following special characters: !@#$%^&*()_+-={}[]|\:,.?
	// The password must contain only printable characters and ascii characters
	// The password must not contain any of the following special characters: '"<>`~;
	passwordValidationResult := strval.ValidateStringWithName(request.NewPassword, "new_password",
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
		return false, passwordValidationResult.Messages
	}

	return true, nil
}

type UserPasswordResetInitiationRequest struct {
	Email string `json:"email"`
}

// Validates the content of a PasswordResetRequest
func (request *UserPasswordResetInitiationRequest) Validate() (valid bool, errors []string) {
	// The email must be not empty and valid email address
	emailValidationResult := strval.ValidateStringWithName(request.Email, "email",
		strval.MustNotBeEmpty(),
		strval.MustBeValidEmailFormat())

	if !emailValidationResult.Valid {
		return false, emailValidationResult.Messages
	}

	return true, nil
}
