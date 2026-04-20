package idam

import "time"

type IdamUserType uint8

const (
	StandardUserType IdamUserType = iota
)

type MembershipStatus uint8

const (
	Active          MembershipStatus = iota // 0 — fully enrolled, has features
	PendingApproval                         // 1 — awaiting admin approval
)

type UserApplication struct {
	// ID of the application that the user belongs to.
	Id string `json:"id"`
	// Slug of the application this user belongs to.
	Slug string `json:"slug"`
	// The features of the application this user has access to
	Features []string `json:"features"`
	// The membership status of the user within this application
	Status MembershipStatus `json:"status"`
}

type User struct {
	Id           string            `json:"id"`
	Username     string            `json:"username"`
	Email        string            `json:"email"`
	Verified     bool              `json:"verified"`
	Type         IdamUserType      `json:"type"`
	Provider     string            `json:"provider"`
	CreatedAtUTC time.Time         `json:"created_at_utc"`
	Applications []UserApplication `json:"applications"`
}
