package appreg

import "testing"

func TestMatchRoute(t *testing.T) {
	tests := []struct {
		requestedURL, pattern string
		expected              bool
	}{
		{"/api/v1/users/123", "/api/v1/users/*", true},
		{"/api/v1/users/123", "/api/v1/users/*/", false},
		{"api/v1/users/123", "/api/v1/users/*", false},
		{"/api/v1/users/123/", "/api/v1/users/*", false},
		{"/api/v1/users/123", "api/v1/users/*", false},
		{"/v1/users/123", "/api/v1/users/*", false},
		{"/api/v1/users/123", "/v1/users/*", false},
		{"", "/v1/users/*", false},
		{"/v1/users/123", "", false},
		{"/api/v1/tenants/123", "/v1/users/*", false},
		{"/api/v1/users/123", "/v1/tenants/*", false},
		{"/api/v1/users/123", "/api/v1/users/*/something/*", false},
		{"/api/v1/users/123/something/456", "/api/v1/users/*", false},
		{"/api/v1/users/123", "/api/v1/users/456", false},
		{"/123/hello", "/*/hello", true},
		{"/123/hello", "*/hello", false},
		{"", "", true},
		{"/", "/", true},
		{"/123", "/*", true},
		{"123", "*", true},
	}

	for _, test := range tests {
		testName := test.pattern + " -> " + test.requestedURL
		t.Run(testName, func(t *testing.T) {
			result := MatchRoute(test.pattern, test.requestedURL)

			if result != test.expected {
				t.Errorf("Expected %t, got %t", test.expected, result)
			}
		})
	}
}

func TestApplyRouteWildcards(t *testing.T) {
	tests := []struct {
		source, dest, expectedValue string
		err                         error
	}{
		// Single-wildcard, length-matched: substitution succeeds.
		{"/api/v1/users/123", "/api/v1/users/*", "/api/v1/users/123", nil},
		// Chained wildcards now succeed when source has matching segment count
		// (previously failed due to token-based resolver looking up "*" literally).
		{"/api/v1/users/123/something/456", "/api/v1/users/*/something/*", "/api/v1/users/123/something/456", nil},
		// Leading wildcard now succeeds when length matches and source segment is non-empty
		// (previously failed because the token-based resolver required a non-empty preceding segment).
		{"/123/hello", "/*/hello", "/123/hello", nil},
		// No wildcards in destination: returned as-is.
		{"/api/v1/users/123", "/api/v1/users/456", "/api/v1/users/456", nil},

		// Empty inputs.
		{"", "/v1/users/*", "", ErrWildcardOperation},
		{"/api/v1/users/123", "", "", ErrWildcardOperation},

		// Length mismatches now error consistently. Several of these previously
		// "succeeded" via the token-based resolver scanning the source URL for
		// the segment preceding "*" — fragile behavior that callers should not
		// have relied on, and that the proxy never exercised because the route
		// matcher rejects mismatched-length pairs upstream.
		{"/api/v1/users/123", "/api/v1/users/*/", "", ErrWildcardOperation},
		{"api/v1/users/123", "/api/v1/users/*", "", ErrWildcardOperation},
		{"/api/v1/users/123/", "/api/v1/users/*", "", ErrWildcardOperation},
		{"/api/v1/users/123", "api/v1/users/*", "", ErrWildcardOperation},
		{"/v1/users/123", "/api/v1/users/*", "", ErrWildcardOperation},
		{"/api/v1/users/123", "/v1/users/*", "", ErrWildcardOperation},
		{"/api/v1/tenants/123", "/v1/users/*", "", ErrWildcardOperation},
		{"/api/v1/users/123", "/v1/tenants/*", "", ErrWildcardOperation},
		{"/api/v1/users/123/something/456", "/api/v1/users/*", "", ErrWildcardOperation},
		{"/123/hello", "*/hello", "", ErrWildcardOperation},
	}

	for _, test := range tests {
		testName := test.source + " -> " + test.dest
		t.Run(testName, func(t *testing.T) {
			result, err := ApplyRouteWildcards(test.source, test.dest)

			if err != test.err {
				t.Errorf("Expected error %v, got %v", test.err, err)
			}

			if result != test.expectedValue {
				t.Errorf("Expected %s, got %s", test.expectedValue, result)
			}
		})
	}
}
