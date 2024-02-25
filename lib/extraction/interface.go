package extraction

import "fmt"

// Common API for extraction

var ErrNoMatches = fmt.Errorf("no data was extracted")

func errNoMatchesf(reason string, args ...any) error {
	reason = fmt.Sprintf(reason, args...)
	return fmt.Errorf("%w (%s)", ErrNoMatches, reason)
}
