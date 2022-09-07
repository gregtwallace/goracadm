package idrac

import (
	"fmt"
)

// ReturnCode is a special string represening "RC" as
// returned in an idrac's XML response.
type ReturnCode string

// known RC meanings
var (
	// success
	RcOK = ReturnCode("0x0")

	// idrac6
	RcIdrac6InvalidUserPassword = ReturnCode("0x10")

	// idrac7
	RcIdrac7InvalidUserPassword = ReturnCode("0x140004")
)

// Error() implements the error interface by returning the
// error code and any known meaning.
func (rc ReturnCode) Error() string {
	return fmt.Sprintf("rc: %s (%s)", string(rc), rc.meaning())
}

// meaning() contains a list of known RC meanings
func (rc *ReturnCode) meaning() string {
	switch *rc {
	case RcOK:
		return "ok"

	case RcIdrac6InvalidUserPassword:
		return "login failed: invalid username or password"
	case RcIdrac7InvalidUserPassword:
		return "login failed: invalid username or password"

	default:
		// break for unknown
	}

	return "meaning unknown"
}
