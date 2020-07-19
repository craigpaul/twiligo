package dates

import (
	"fmt"
	"strings"
	"time"
)

// Rfc2822Time ...
type Rfc2822Time struct {
	time.Time
}

// UnmarshalJSON ...
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

// MarshalJSON ...
func (t Rfc2822Time) MarshalJSON() ([]byte, error) {
	return []byte(t.String()), nil
}

func (t *Rfc2822Time) String() string {
	format := "Mon, 02 Jan 2006 15:04:05 -0700"

	return fmt.Sprintf("%q", t.Format(format))
}
