package labs

import (
	"time"

	"github.com/dmars8047/marshall-labs-sdk/idam"
)

// The event types
type EventType uint8

// The event type values
const (
	// The event type for when a user is created
	UserCreatedEventType EventType = iota
	// The event type for when a user is updated
	UserUpdatedEventType
	// The event type for when a user is deleted
	UserDeletedEventType
	// The event type for when a user logs out
	UserLoggedOutEventType
	// The event type for when a user becomes an active member of an application
	MemberActivatedEventType
	// The event type for when a user ceases to be an active member of an application
	MemberDeactivatedEventType
)

// Event action names. These form the middle segment of the
// "user-action.<action>.<appId>" RabbitMQ routing key published by idam-service.
const (
	// A user became an active member of an application (registration verified,
	// invite accepted, or admin-approved). Subsumes the former "account-verified".
	MemberActivatedAction = "member-activated"
	// A user lost active membership of an application (account deletion, leave, kick).
	MemberDeactivatedAction = "member-deactivated"
	// A user account was deleted entirely.
	UserDeletedAction = "user-deleted"
)

// Contains common idam event properties
type Event struct {
	// The type identifier of the event
	EventType EventType `json:"event_type"`
	// The UTC timestamp of the event
	TimeStamp time.Time `json:"time_stamp"`
}

type UserLoggedOutEvent struct {
	Event
	// The ID of the user that logged out
	UserId string `json:"user_id"`
	// The ID of the token that was invalidated as a result of the user log out
	InvalidatedTokenId string `json:"invalidated_token_id"`
	// The ID of the application the log out event occurred under
	ApplicationId string `json:"application_id"`
}

type UserCreatedEvent struct {
	Event
	// The ID of the user that was created in the IDAM system
	User idam.User `json:"user"`
}

// MemberActivatedEvent is published when a user becomes an active member of an
// application. It carries enough for a consumer to provision a profile on first
// sight (Email) and record the membership association. The activation timestamp
// is the embedded Event.TimeStamp.
type MemberActivatedEvent struct {
	Event
	// The ID of the user that became an active member.
	UserId string `json:"user_id"`
	// The ID of the application the user is now an active member of.
	ApplicationId string `json:"application_id"`
	// The user's email, used to seed a display name when the profile is first created.
	Email string `json:"email"`
}

// MemberDeactivatedEvent is published when a user ceases to be an active member
// of an application. The embedded Event.TimeStamp is the deactivation time.
type MemberDeactivatedEvent struct {
	Event
	// The ID of the user that was deactivated.
	UserId string `json:"user_id"`
	// The ID of the application the user is no longer an active member of.
	ApplicationId string `json:"application_id"`
}

// UserDeletedEvent is published when a user account is removed entirely.
type UserDeletedEvent struct {
	Event
	// The ID of the user that was deleted.
	UserId string `json:"user_id"`
}
