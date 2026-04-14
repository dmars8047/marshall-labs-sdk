package labs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/dmars8047/marshall-labs-sdk/appreg"
	"github.com/dmars8047/marshall-labs-sdk/idam"
)

// IdamSecurityClient is a client which facilitates communication with the IDAM service's security endpoints
type IdamSecurityClient struct {
	httpClient *http.Client
	baseUrl    string
	hmacKey    []byte
}

// Creates a new IdamSecurityClient.
// hmacKey must match the key configured in the idam service's UserAuthGuard middleware.
func NewIdamSecurityClient(httpClient *http.Client, baseUrl string, hmacKey []byte) *IdamSecurityClient {
	return &IdamSecurityClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
		hmacKey:    hmacKey,
	}
}

const (
	VerifyTokenEndpointUrlSuffix = "/api/idam/security/verify-token"
	VerifyRouteEndpointUrlSuffix = "/api/idam/security/verify-route"
	GetPublicRoutesUrl           = "/api/idam/appreg/public-routes"
)

// Function to authenticate the JWT token using the IDAM service
func (idamAuthService *IdamSecurityClient) VerifyToken(token string) (*VerifyTokenResponse, error) {

	// send a post request with an empty body, make sure the authentication header is attached
	// if the response is 200, return nil
	url, err := url.Parse(idamAuthService.baseUrl + VerifyTokenEndpointUrlSuffix)

	if err != nil {
		return nil, err
	}

	// Create a new request to forward to the target
	verifyRequest, err := http.NewRequest(http.MethodPost, url.String(), nil)

	if err != nil {
		return nil, err
	}

	// Add the authorization header to the request
	verifyRequest.Header.Add("Authorization", token)

	resp, err := idamAuthService.httpClient.Do(verifyRequest)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// UnMarhsal the response body into an ErrorResponse object
		var errorResponse idam.ErrorResponse

		err = json.NewDecoder(resp.Body).Decode(&errorResponse)

		if err != nil {
			return nil, fmt.Errorf("error decoding response body from idam service - %v", err)
		}

		return nil, &errorResponse
	}

	var verifyTokenResponse VerifyTokenResponse

	err = json.NewDecoder(resp.Body).Decode(&verifyTokenResponse)

	if err != nil {
		return nil, errors.New("error decoding response body from idam service")
	}

	return &verifyTokenResponse, nil
}

// Function to authenticate the JWT token using the IDAM service
func (idamAuthService *IdamSecurityClient) VerifyRoute(requestedRoute, requestedVerb string, identityContext *UserAuthContext) (*appreg.Route, error) {

	// Call the IDAM service /api/idam/auth/verify endpoint and validate the token via http
	url, err := url.Parse(idamAuthService.baseUrl + VerifyRouteEndpointUrlSuffix)

	if err != nil {
		return nil, err
	}

	verifyRequestObj := VerifyRouteRequest{
		RequestedRoute: requestedRoute,
		RequestedVerb:  requestedVerb,
	}

	requestBodyBytes, err := json.Marshal(verifyRequestObj)

	if err != nil {
		return nil, err
	}

	verifyRequest, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(requestBodyBytes))

	if err != nil {
		return nil, err
	}

	identityContext.ApplyRequestHeaders(verifyRequest, idamAuthService.hmacKey)

	// Add the content type header to the request
	verifyRequest.Header.Add("Content-Type", "application/json")

	response, err := idamAuthService.httpClient.Do(verifyRequest)

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

	var route appreg.Route

	err = json.NewDecoder(response.Body).Decode(&route)

	if err != nil {
		return nil, errors.New("error decoding response body from idam service")
	}

	return &route, nil
}
