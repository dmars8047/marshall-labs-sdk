package appreg

import (
	"errors"
	"strings"
)

var ErrWildcardOperation = errors.New("invalid wildcard matching operation")

// MatchRoute reports whether requestedURL matches a route pattern. The pattern
// may contain "*" segments that each match any single path segment. Segment
// counts must be equal — "*" never spans multiple segments.
func MatchRoute(pattern, requestedURL string) bool {
	patternParts := strings.Split(pattern, "/")
	requestedParts := strings.Split(requestedURL, "/")

	if len(patternParts) != len(requestedParts) {
		return false
	}

	for i := range patternParts {
		if patternParts[i] == "*" {
			continue
		}

		if patternParts[i] != requestedParts[i] {
			return false
		}
	}

	return true
}

// ApplyRouteWildcards substitutes "*" segments in destinationPattern with the
// segment at the same index from sourceURL. Callers must have already verified
// that the two URLs match via MatchRoute (segment counts equal).
func ApplyRouteWildcards(sourceURL, destinationPattern string) (string, error) {
	if len(sourceURL) < 1 || len(destinationPattern) < 1 {
		return "", ErrWildcardOperation
	}

	if !strings.Contains(destinationPattern, "*") {
		return destinationPattern, nil
	}

	sourceParts := strings.Split(sourceURL, "/")
	destinationParts := strings.Split(destinationPattern, "/")

	if len(sourceParts) != len(destinationParts) {
		return "", ErrWildcardOperation
	}

	for i, part := range destinationParts {
		if part != "*" {
			continue
		}

		if sourceParts[i] == "" {
			return "", ErrWildcardOperation
		}

		destinationParts[i] = sourceParts[i]
	}

	return strings.Join(destinationParts, "/"), nil
}
