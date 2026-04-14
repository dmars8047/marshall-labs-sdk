package idam

// ErrorResponse is the response returned when an error
// is encoutered during the processing of a request to the IDAM API
type ErrorResponse struct {
	Code    uint16   `json:"error_code"`
	Message string   `json:"error_message"`
	Details []string `json:"error_details"`
	// Indicates the marshall-labs application (slug) that the error originated from.
	SourceApp string `json:"error_source_app"`
}

// Error returns the error message for the ErrorResponse
func (err ErrorResponse) Error() string {
	return err.Message
}

// NewErrorResponse creates an ErrorResponse with the given code and message.
// Usage: NewErrorResponse(0, "An error occured during validation")
func NewErrorResponse(code uint16, message string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   []string{},
		SourceApp: "idam",
	}
}

// NewDetailedErrorResponse creates an ErrorResponse with the given code, message and details.
// Usage: NewDetailedErrorResponse(0, "An error occured during validation","message1", "message2")
func NewDetailedErrorResponse(code uint16, message string, details ...string) *ErrorResponse {
	return &ErrorResponse{
		Code:      code,
		Message:   message,
		Details:   details,
		SourceApp: "idam",
	}
}

// NewUnhandledErrorResponse creates an ErrorResponse with the code and message for an unhandled error.
func NewUnhandledErrorResponse() *ErrorResponse {
	return &ErrorResponse{
		Code:      UnhandledError,
		Message:   UnhandledErrorMessage,
		Details:   []string{},
		SourceApp: "idam",
	}
}

const (
	// Error codes
	// Error code 0 indicates an unhandled error. This means there was a server error.
	UnhandledError        = 1
	UnhandledErrorMessage = "an unhandled/unexpected error occured"
	// Error code 5 indicates the request body could not be parsed or was otherwise invalid.
	RequestPayloadInvalid     = 5
	RequestBodyInvalidMessage = "the request body could not be parsed"
	// Error code 10 indicates the request failed validation.
	// This means the request content was parsed but failed validation of the content.
	RequestValidationFailure        = 10
	RequestValidationFailureMessage = "request validation failure"
	// Error code 15 indicates the requested application resource was not found.
	ApplicationNotFound        = 15
	ApplicationNotFoundMessage = "application not found"
	// Error code 20 indicates the credentials provided were invalid.
	InvalidCredentials        = 20
	InvalidCredentialsMessage = "invalid credentials"
	// Error code 25 indicates the data provided conflicts with existing data.
	// This means that the data provided cannot be used because it conflicts with existing data.
	DataConflict        = 25
	DataConflictMessage = "data conflict"
	// Error code 30 indicates the user has not verified their email address.
	UserNotVerified        = 30
	UserNotVerifiedMessage = "user not verified"
	// Error code 35 indicates the provided authorization token was invalid or has been blacklisted.
	InvalidAuthToken        = 35
	InvalidAuthTokenMessage = "invalid or malformed authorization token"
	// Error code 40 indicates the user does not have access to the requested resource.
	AccessDenied        = 40
	AccessDeniedMessage = "access denied"
	// Error code 45 indicates the provided verification code was invalid.
	InvalidUserVerficationToken        = 45
	InvalidUserVerficationTokenMessage = "invalid verification code"
	// Error code 50 indicates the requested user resource was not found.
	UserNotFound        = 50
	UserNotFoundMessage = "user not found"
	// Error code 55 indicates the provided password reset token was invalid.
	InvalidPasswordResetToken        = 55
	InvalidPasswordResetTokenMessage = "invalid password reset token"
	// Error code 60 indicates the provided password reset verification code was invalid.
	InvalidPasswordResetVerificationCode        = 60
	InvalidPasswordResetVerificationCodeMessage = "invalid password reset verification code"
	// Error code 65 indicates that missing or invalid request headers were provided.
	InvalidRequestHeaders        = 65
	InvalidRequestHeadersMessage = "invalid or missing request headers"
	// Error code 70 indicates that an authorization token is expired.
	AuthTokenExpired        = 70
	AuthTokenExpiredMessage = "authorization token expired"
	// Error code 75 indicates the user's account is locked out due to too many failed login attempts.
	// The lockout will expire 1 hour from the user's last failed login attempt.
	UserAccountLockout        = 75
	UserAccountLockoutMessage = "user account lockout due to too many failed login attempts"
	// Error code 80 indicates that the user is not registered with the target application
	UserNotRegisteredWithApplication        = 80
	UserNotRegisteredWithApplicationMessage = "user not registered with target application"
	// Error code 85 is similar to a Data Conflict but is specific to a users email
	UserEmailConflict        = 85
	UserEmailConflictMessage = "user email conflict"
	// Error code 90 is similar to a Data Conflict but is specific to a users username
	UsernameConflict        = 95
	UsernameConflictMessage = "username conflict"
)
