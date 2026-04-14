package labs

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dmars8047/marshall-labs-sdk/appreg"
)

// IdamPublicRoutesClient fetches the public route list from the idam-service.
type IdamPublicRoutesClient struct {
	httpClient *http.Client
	baseUrl    string
}

func NewIdamPublicRoutesClient(httpClient *http.Client, baseUrl string) *IdamPublicRoutesClient {
	return &IdamPublicRoutesClient{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (c *IdamPublicRoutesClient) GetPublicRoutes() ([]appreg.Route, error) {
	resp, err := c.httpClient.Get(c.baseUrl + GetPublicRoutesUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status fetching public routes from idam service")
	}

	var routes []appreg.Route
	if err = json.NewDecoder(resp.Body).Decode(&routes); err != nil {
		return nil, errors.New("error decoding public routes response from idam service")
	}

	return routes, nil
}
