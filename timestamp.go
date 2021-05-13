package null

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Timestamp is a nullable time.Time. It supports SQL and JSON serialization.
// It will marshal to null if null.
type Timestamp struct {
	sql.NullTime
}

// Value implements the driver Valuer interface.
func (t Timestamp) Value() (driver.Value, error) {
	if !t.Valid {
		return nil, nil
	}
	return t.Time, nil
}

// NewTimestamp creates a new Timestamp.
func NewTimestamp(t time.Time, valid bool) Timestamp {
	return Timestamp{
		NullTime: sql.NullTime{
			Time:  t,
			Valid: valid,
		},
	}
}

// TimestampFrom creates a new Timestamp that will always be valid.
func TimestampFrom(t time.Time) Timestamp {
	return NewTimestamp(t, true)
}

// TimestampFromPtr creates a new Timestamp that will be null if t is nil.
func TimestampFromPtr(t *time.Time) Timestamp {
	if t == nil {
		return NewTimestamp(time.Time{}, false)
	}
	return NewTimestamp(*t, true)
}

// ValueOrZero returns the inner value if valid, otherwise zero.
func (t Timestamp) ValueOrZero() time.Time {
	if !t.Valid {
		return time.Time{}
	}
	return t.Time
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this timestamp is null.
func (t Timestamp) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports int64 and null input.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		t.Valid = false
		return nil
	}
	var v int64
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("null: couldn't unmarshal JSON: %w", err)
	}
	t.Time = time.Unix(v, 0)
	t.Valid = true
	return nil
}

// MarshalText implements encoding.TextMarshaler.
// It returns an empty string if invalid, otherwise int64.
func (t Timestamp) MarshalText() ([]byte, error) {
	if !t.Valid {
		return []byte{}, nil
	}
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
// It will unmarshal to a null int64 Unix timestamp to time.Time if the input is a blank or not an time.Time.
func (t *Timestamp) UnmarshalText(text []byte) error {
	str := string(text)
	// allowing "null" is for backwards compatibility with v3
	if str == "" || str == "null" {
		t.Valid = false
		return nil
	}
	v, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return fmt.Errorf("null: couldn't unmarshal text: %w", err)
	}
	t.Time = time.Unix(v, 0)
	t.Valid = true
	return nil
}

// SetValid changes this Timestamp's value and sets it to be non-null.
func (t *Timestamp) SetValid(v time.Time) {
	t.Time = v
	t.Valid = true
}

// Ptr returns a pointer to this Timestamp's value, or a nil pointer if this Time is null.
func (t Timestamp) Ptr() *time.Time {
	if !t.Valid {
		return nil
	}
	return &t.Time
}

// IsZero returns true for invalid Times, hopefully for future omitempty support.
// A non-null Time with a zero value will not be considered zero.
func (t Timestamp) IsZero() bool {
	return !t.Valid
}

// Equal returns true if both Timestamp objects encode the same time or are both null.
// Two times can be equal even if they are in different locations.
// For example, 6:00 +0200 CEST and 4:00 UTC are Equal.
func (t Timestamp) Equal(other Timestamp) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time.Equal(other.Time))
}

// ExactEqual returns true if both Timestamp objects are equal or both null.
// ExactEqual returns false for times that are in different locations or
// have a different monotonic clock reading.
func (t Timestamp) ExactEqual(other Timestamp) bool {
	return t.Valid == other.Valid && (!t.Valid || t.Time == other.Time)
}
