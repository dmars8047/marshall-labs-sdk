package idam

import "testing"

func TestUserPasswordResetExecutionRequest_Validate_Password(t *testing.T) {
	validBase := UserPasswordResetExecutionRequest{
		UserID:             "user-id",
		PasswordResetToken: "reset-token",
		VerificationCode:   "verify-code",
	}

	for _, tt := range passwordValidationCases {
		t.Run(tt.name, func(t *testing.T) {
			req := validBase
			req.NewPassword = tt.password
			valid, errs := req.Validate()
			if valid != tt.wantValid {
				t.Errorf("new_password %q: got valid=%v, want %v; errors: %v", tt.password, valid, tt.wantValid, errs)
			}
		})
	}
}
