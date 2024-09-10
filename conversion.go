package xrpl

import (
	"time"
)

/*
 * Ripple time to Unix time conversion
 */

// Seconds since UNIX Epoch to Ripple Epoch (2000-01-01T00:00 UTC)
// https://xrpl.org/docs/references/protocol/data-types/basic-data-types#specifying-time
const RIPPLE_EPOCH_DIFF int64 = 946684800

// Convert a Ripple timestamp to a unix timestamp.
func RippleTimeToUnixTime(rippleTime int64) int64 {
	return (rippleTime + RIPPLE_EPOCH_DIFF)
}

// Convert a unix timestamp to a Ripple timestamp.
func UnixTimeToRippleTime(unixTime int64) int64 {
	return (unixTime - RIPPLE_EPOCH_DIFF)
}

// Convert a Ripple timestamp to an ISO8601 time.
func RippleTimeToISOTime(rippleTime int64) string {
	unixTime := RippleTimeToUnixTime(rippleTime)
	return time.Unix(unixTime, 0).UTC().Format(time.RFC3339)
}

// Convert an ISO8601 timestamp to a Ripple timestamp.
func IsoTimeToRippleTime(isoTime string) (int64, error) {
	theTime, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return 0, err
	}
	rippleTime := UnixTimeToRippleTime(theTime.Unix())
	return rippleTime, nil
}
