package labs

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/dmars8047/marshall-labs-sdk/idam"
)

const (
	GetUserUrl = "/api/idam/users/:usrId"
	// Internal, service-auth-guarded, Docker-network-only endpoints (not routed via api-proxy).
	GetAppMembersUrl   = "/api/idam/internal/applications/:appId/members"
	GetInternalUserUrl = "/api/idam/internal/users/:userId"

	requestIDHeader = "X-Request-ID"
)

// A client for making internal http calls to the IDAM service's user centered endpoints.
// This is distinct from the user facing user account endpoint client.
// This client should be used for making internal calls.
type IdamUserClient struct {
	httpClient  *http.Client
	baseUrl     string
	serviceName string
	signingKey  []byte
}

// NewIdamUserClient builds a client for idam's internal user endpoints. serviceName and
// signingKey are used to sign the X-Service-Auth header on the internal, service-guarded
// calls (GetAppMembers, GetInternalUser); signingKey must match idam's IDAM_HEADER_SIGNING_KEY.
func NewIdamUserClient(httpClient *http.Client, baseUrl, serviceName string, signingKey []byte) *IdamUserClient {
	return &IdamUserClient{
		httpClient:  httpClient,
		baseUrl:     baseUrl,
		serviceName: serviceName,
		signingKey:  signingKey,
	}
}

// authedGet issues a GET to the resolved path with a signed X-Service-Auth header and,
// when traceID is non-empty, an X-Request-ID header so the call is traceable across the hop.
func (client *IdamUserClient) authedGet(path, traceID string) (*http.Response, error) {
	base, err := url.Parse(client.baseUrl)
	if err != nil {
		return nil, err
	}

	suffix, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	resolved := base.ResolveReference(suffix)

	req, err := http.NewRequest(http.MethodGet, resolved.String(), nil)
	if err != nil {
		return nil, err
	}

	authHeader, err := CreateServiceAuthHeader(client.serviceName, client.signingKey)
	if err != nil {
		return nil, err
	}
	req.Header.Set(ServiceAuthHeader, authHeader)

	if traceID != "" {
		req.Header.Set(requestIDHeader, traceID)
	}

	return client.httpClient.Do(req)
}

// decodeIdamError decodes an idam ErrorResponse from a non-200 response body.
func decodeIdamError(response *http.Response) error {
	var errorResponse idam.ErrorResponse
	if err := json.NewDecoder(response.Body).Decode(&errorResponse); err != nil {
		return fmt.Errorf("error decoding response body from idam service - %v", err)
	}
	return &errorResponse
}

func (client *IdamUserClient) GetUser(userId string) (*idam.User, error) {
	urlSuffix := strings.Replace(GetUserUrl, ":usrId", userId, 1)

	// Parse the base URL
	base, err := url.Parse(client.baseUrl)

	if err != nil {
		return nil, err
	}

	// Parse the suffix as a URL
	suffix, err := url.Parse(urlSuffix)

	if err != nil {
		return nil, err
	}

	// Resolve to correct URL
	resolvedURL := base.ResolveReference(suffix)

	response, err := client.httpClient.Get(resolvedURL.String())

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, decodeIdamError(response)
	}

	var user idam.User

	err = json.NewDecoder(response.Body).Decode(&user)

	if err != nil {
		return nil, errors.New("error decoding register user response body from idam service")
	}

	return &user, nil
}

// GetInternalUser fetches a single idam user via the internal, service-auth-guarded endpoint.
// Used by downstream services (e.g. user-service) to look up a user without a user-auth context.
func (client *IdamUserClient) GetInternalUser(userID, traceID string) (*idam.User, error) {
	path := strings.Replace(GetInternalUserUrl, ":userId", userID, 1)

	response, err := client.authedGet(path, traceID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, decodeIdamError(response)
	}

	var user idam.User
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return nil, errors.New("error decoding user response body from idam service")
	}

	return &user, nil
}

// GetAppMembers returns the active members of an application, each with the email,
// features, status, and join time idam holds for them. It targets idam's internal
// service-auth-guarded endpoint and is used to compose the platform member directory.
func (client *IdamUserClient) GetAppMembers(appID, traceID string) ([]idam.AppMember, error) {
	path := strings.Replace(GetAppMembersUrl, ":appId", appID, 1)

	response, err := client.authedGet(path, traceID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, decodeIdamError(response)
	}

	var members []idam.AppMember
	if err := json.NewDecoder(response.Body).Decode(&members); err != nil {
		return nil, errors.New("error decoding app members response body from idam service")
	}

	return members, nil
}
