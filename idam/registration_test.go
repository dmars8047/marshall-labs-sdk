package idam

import "testing"

func TestUserRegistrationRequest_Validate_Password(t *testing.T) {
	validBase := UserRegistrationRequest{
		Username: "testuser",
		Email:    "test@example.com",
	}

	for _, tt := range passwordValidationCases {
		t.Run(tt.name, func(t *testing.T) {
			req := validBase
			req.Password = tt.password
			valid, errs := req.Validate()
			if valid != tt.wantValid {
				t.Errorf("password %q: got valid=%v, want %v; errors: %v", tt.password, valid, tt.wantValid, errs)
			}
		})
	}
}

func TestUserRegistrationRequest_Validate_AllFieldsValid(t *testing.T) {
	req := UserRegistrationRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "correct horse battery staple",
	}
	valid, errs := req.Validate()
	if !valid {
		t.Errorf("expected valid registration request, got errors: %v", errs)
	}
}
