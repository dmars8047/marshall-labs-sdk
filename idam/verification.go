package idam

type UserAccountVerificationRequest struct {
	UserId            string `json:"user_id"`
	VerificationToken string `json:"verification_token"`
}
