package appreg

import "errors"

type Application struct {
	ID                     string    `json:"id"`
	Name                   string    `json:"name"`
	Description            string    `json:"description"`
	Slug                   string    `json:"slug"`
	Features               []Feature `json:"features"`
	PublicBaseUrl          string    `json:"public_base_url"`
	VerifyUserUrlSuffix    string    `json:"verify_user_url_suffix"`
	PasswordResetUrlSuffix string    `json:"password_reset_url_suffix"`
}

type Feature struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsDefault   bool    `json:"is_default"`
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
