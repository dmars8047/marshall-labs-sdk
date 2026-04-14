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
)

// A client for making internal http calls to the IDAM service's user centered endpoints.
// This is distinct from the user facing user account endpoint client.
// This client should be used for making internal calls.
type IdamUserClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewIdamUserClient(httpClient *http.Client, baseUrl string) *IdamUserClient {
	return &IdamUserClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
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
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse idam.ErrorResponse

		err = json.NewDecoder(response.Body).Decode(&errorResponse)

		if err != nil {
			return nil, fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return nil, &errorResponse
	}

	var user idam.User

	err = json.NewDecoder(response.Body).Decode(&user)

	if err != nil {
		return nil, errors.New("error decoding register user response body from idam service")
	}

	return &user, nil
}
