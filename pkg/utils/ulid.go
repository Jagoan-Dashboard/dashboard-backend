package utils

import (
	"crypto/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

// GenerateULID generates a new ULID string
func GenerateULID() string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}

// ParseULID parses a ULID string and returns the ULID object
func ParseULID(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}

// IsValidULID checks if a string is a valid ULID
func IsValidULID(s string) bool {
	_, err := ulid.Parse(s)
	return err == nil
}

// ULIDFromTime generates a ULID from a specific time
func ULIDFromTime(t time.Time) string {
	entropy := ulid.Monotonic(rand.Reader, 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}