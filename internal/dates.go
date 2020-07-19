package dates

import (
	"fmt"
	"strings"
	"time"
)

// Rfc2822Time is meant to allow access to time.Time but successfully parse the RFC 2822 date format into a struct easily.
type Rfc2822Time struct {
	time.Time
}

// UnmarshalJSON handles converting a string that is hopefully in the RFC 2822 date format into a usable struct with access to the time.Time functions.
func (t *Rfc2822Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	format := "Mon, 02 Jan 2006 15:04:05 -0700"
	nt, err := time.Parse(format, s)

	if err != nil {
		return err
	}

	*t = Rfc2822Time{nt}

	return nil
}

// MarshalJSON handles converting a Rfc2822Time struct into the proper string representation when being serialized to JSON.
func (t Rfc2822Time) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Rfc2822Time) String() string {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"

	return fmt.Sprintf("%q", t.Format(format))
}
