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
	// UTC timestamp of when the user became an active member of this application.
	// Zero for memberships that predate join-time tracking.
	JoinedAtUTC time.Time `json:"joined_at"`
}

// AppMember is an active member of an application, as returned by idam's internal
// directory endpoint. It carries enough for a consumer (e.g. user-service) to build a
// member listing: identity, the email used to seed a display name, the member's
// features within the app, and when they joined.
type AppMember struct {
	UserId      string           `json:"user_id"`
	Email       string           `json:"email"`
	Features    []string         `json:"features"`
	Status      MembershipStatus `json:"status"`
	JoinedAtUTC time.Time        `json:"joined_at"`
}

type User struct {
	Id           string            `json:"id"`
	Email        string            `json:"email"`
	Verified     bool              `json:"verified"`
	Type         IdamUserType      `json:"type"`
	Provider     string            `json:"provider"`
	CreatedAtUTC time.Time         `json:"created_at_utc"`
	Applications []UserApplication `json:"applications"`
}
