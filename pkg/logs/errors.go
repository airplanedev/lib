package logs

import "regexp"

// https://regex101.com/r/dV2iG1/2
var nodeErrorsRegex = regexp.MustCompile(`^Error \[ERR_REQUIRE_ESM\](?:.*\/node_modules\/([^\/]+)\/)?`)

func hasNodeESMError(log string) (string, bool) {
	matches := nodeErrorsRegex.FindStringSubmatch(log)
	if len(matches) == 1 {
		// Matched the error, but could not extract a module name.
		return "", true
	}
	if len(matches) == 2 {
		// Matched the error with a specific module.
		return matches[1], true
	}

	// Not a match.
	return "", false
}
