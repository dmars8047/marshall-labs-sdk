package idam

import "testing"

func TestUserPasswordResetExecutionRequest_Validate_Password(t *testing.T) {
	validBase := UserPasswordResetExecutionRequest{
		UserID:             "user-id",
		PasswordResetToken: "reset-token",
		VerificationCode:   "verify-code",
	}

	tests := []struct {
		name      string
		password  string
		wantValid bool
	}{
		{"too short (11 chars)", "abcdefghijk", false},
		{"minimum length (12 chars, no complexity)", "abcdefghijkl", true},
		{"long passphrase no complexity", "correct horse battery staple", true},
		{"too long (65 chars)", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},
		{"non-ASCII character", "abcdefghijkl\xe9", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
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
