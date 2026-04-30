package appreg

import "errors"

type MembershipType uint8

const (
	OpenEnrollment  MembershipType = iota // 0 — users self-register freely (default)
	AdminApproval                         // 1 — new registrations require admin approval
	InviteOnly                            // 2 — users may only join via admin-issued invitation
)

type Application struct {
	ID                     string         `json:"id"`
	Name                   string         `json:"name"`
	Description            string         `json:"description"`
	Slug                   string         `json:"slug"`
	MembershipType         MembershipType `json:"membership_type"`
	Features               []Feature      `json:"features"`
	PublicBaseUrl          string         `json:"public_base_url"`
	VerifyUserUrlSuffix    string         `json:"verify_user_url_suffix"`
	PasswordResetUrlSuffix string         `json:"password_reset_url_suffix"`
	InviteUrlSuffix        string         `json:"invite_url_suffix"`
	DefaultRedirectUri     string         `json:"default_redirect_uri"`
}

type Feature struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsDefault   bool    `json:"is_default"`
	IsAdmin     bool    `json:"is_admin"`
	Routes      []Route `json:"routes"`
}

type Route struct {
	ID          string `json:"id"`
	UrlSuffix   string `json:"url_suffix"`
	ServiceUrl  string `json:"service_url"`
	ServiceName string `json:"service_name"`
	Port        uint16 `json:"port"`
	Verb        string `json:"verb"`
	Public      bool   `json:"is_public"`
}

type Service struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	HostUrl     string  `json:"host_url"`
	Port        uint16  `json:"port"`
	Routes      []Route `json:"routes"`
}

var ErrApplicationNotFound = errors.New("application not found")
var ErrServiceNotFound = errors.New("service not found")
