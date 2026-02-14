package maxigo

import "encoding/json"

// Optional represents a value that may or may not be set.
// Use [Some] to create a set value. The zero value is unset.
//
// When used with JSON struct tag "omitzero", unset fields are omitted
// from marshaled output. This allows distinguishing between three states:
//   - Unset (zero value): field omitted from JSON
//   - Set to zero value: field included (e.g. "" for string, false for bool)
//   - Set to non-zero value: field included with value
type Optional[T any] struct {
	Value T
	Set   bool
}

// Some creates a set [Optional] with the given value.
func Some[T any](v T) Optional[T] {
	return Optional[T]{Value: v, Set: true}
}

// IsZero reports whether the Optional is unset.
// Used by encoding/json with the "omitzero" tag option.
func (o Optional[T]) IsZero() bool {
	return !o.Set
}

// MarshalJSON marshals the underlying value.
// If the Optional is unset, it marshals as JSON null.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.Set {
		return []byte("null"), nil
	}
	return json.Marshal(o.Value)
}

// UnmarshalJSON unmarshals the value and marks the Optional as set.
// JSON null is treated as unset (Set remains false).
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	o.Set = true
	return json.Unmarshal(data, &o.Value)
}

// Common type aliases for convenience.
type (
	OptString = Optional[string]
	OptBool   = Optional[bool]
	OptInt64  = Optional[int64]
)