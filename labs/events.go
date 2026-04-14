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
