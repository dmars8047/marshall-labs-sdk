package idam

import "time"

type InviteStatus string

const (
	InviteStatusPending  InviteStatus = "pending"
	InviteStatusAccepted InviteStatus = "accepted"
	InviteStatusRevoked  InviteStatus = "revoked"
	InviteStatusExpired  InviteStatus = "expired"
)

type SendInviteRequest struct {
	InviteeEmail string `json:"invitee_email"`
}

type AcceptInviteRequest struct {
	InviteToken string `json:"invite_token"`
}

type InviteListItem struct {
	Id           string       `json:"id"`
	InviteeEmail string       `json:"invitee_email"`
	InvitedBy    string       `json:"invited_by"`
	CreatedAt    time.Time    `json:"created_at"`
	ExpiresAt    time.Time    `json:"expires_at"`
	AcceptedAt   *time.Time   `json:"accepted_at,omitempty"`
	RevokedAt    *time.Time   `json:"revoked_at,omitempty"`
	Status       InviteStatus `json:"status"`
}
