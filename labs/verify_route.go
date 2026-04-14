package labs

// VerifyRouteRequest is the request body for the VerifyRoute endpoint
type VerifyRouteRequest struct {
	// The route that is being requested
	RequestedRoute string `json:"requested_route"`
	// The verb that is being requested
	RequestedVerb string `json:"requested_verb"`
}
