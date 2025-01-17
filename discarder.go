package wrappers

import (
	"encoding/json"
	"reflect"
)

// Discarder is a generic type that proxies any WrapperProvider.
// It handles  unmarshalling by suppressing errors of discarded invalid values.
type Discarder[W WrapperProvider] struct {
	Proxy W
}

// MarshalJSON marshals the underlying wrapper to JSON.
func (discarder *Discarder[W]) MarshalJSON() ([]byte, error) {
	if discarder.Proxy.IsDiscarded() {
		return json.Marshal(nil)
	}
	return discarder.Proxy.MarshalJSON()
}

// UnmarshalJSON unmarshals JSON data into the underlying wrapper.
// If unmarshalling fails, it discards the value without returning an error.
func (discarder *Discarder[W]) UnmarshalJSON(data []byte) error {
	err := discarder.Proxy.UnmarshalJSON(data)
	if err != nil {
		// Check if our proxy is nil via reflection.
		if reflect.ValueOf(discarder.Proxy).IsNil() {
			return nil
		}

		discarder.Proxy.Discard()
		return nil // Suppress the error by Discarder
	}
	return nil
}

// NewDiscarder initializes a new Discarder wrapper with the provided underlying wrapper.
func NewDiscarder[W WrapperProvider](wrapper W) *Discarder[W] {
	Discarder := &Discarder[W]{
		Proxy: wrapper,
	}
	return Discarder
}
