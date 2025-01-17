package wrappers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	WrappersTagHeader  = "wrappers"
	WrappersTagDiscard = "discard"
)

// Wrapper is an interface that wraps a value. It has two methods: Wrap and Unwrap.

// Wrappers are used to perform type assertions and validations on values before assigment and later retrieval.
// Their main value add is that they commonly implement the Wrapper interface, which allows them to be used in a generic way.

// When wrapping a value, the Wrap method is called with the value to be wrapped and a boolean indicating if the value should be discarded if it is invalid.

// Name is a type that holds the name of the wrapper. It is used to identify the wrapper in error messages.
type Name string

// WrapperBase is a struct that holds the basic fields of a wrapper. It is embedded in all wrapper implementations.
type WrapperBase struct {
	initialized bool // Indicates if the wrapper has been initialized.
	discarded   bool // If this is true, unwrapping will return nil. This is useful when we want to discard for processes where we need to explicitly exclude data such as during an API call where we shouldn't send a field.
}

func (wrapper *WrapperBase) Initialize() {
	wrapper.initialized = true
}

func (wrapper *WrapperBase) IsInitialized() bool {
	return wrapper.initialized
}

func (wrapper *WrapperBase) Discard() {
	wrapper.discarded = true
}

func (wrapper *WrapperBase) IsDiscarded() bool {
	return wrapper.discarded
}

// WrapperProvider is an interface that defines the methods that a wrapper must implement.
type WrapperProvider interface {
	// The initization methods are important in cases where parameters or other custom logic is needed before the wrapper can be used.
	// An example of this would be derivitives of the WrapperRegex wrapper where we need to set the regex pattern before we can use the wrapper.
	Initialize()         // Initializes the wrapper.
	IsInitialized() bool // Returns true if the wrapper has been initialized.

	Discard()          // Discards the value. Sets the Discard flag to true.
	IsDiscarded() bool // Returns true if the value was discarded. This method is important as it is called during marshalling to JSON. If the value was nullified, we should return nil.

	Wrap(any, bool) error // Wraps a value. The value is validated and stored in the wrapper with the wrappers type. The discard parameter indicates if the value should be discarded if it is invalid without returning an error.

	// We need to implement the MarshalJSON and UnmarshalJSON methods in order to be able to use the wrappers in JSON marshalling and unmarshalling.
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error

	UnwrapAny() any // Similar to Unwrap, but returns the value as an interface{}.
	GetAny() any    // Similar to Get, but returns the value as an interface{}.
}

type UnwrapResult interface {
	bool | int64 | float64 | string // The unwrap method always returns a value of this type.
}

// The core wrapper struct used for all implementations. Importantly, it is a generic implementation but embeds the WrapperBase struct.
type Wrapper[V any, R any] struct {
	WrapperBase
	Value V // The value that is wrapped.
}

// The main implementation of a Wrapper. This is the core implementation that all other wrappers should implement.
type WrapperImplementation[V any, R UnwrapResult] interface {
	WrapperProvider

	Unwrap() R // Gets the stored value as one of the types specified in the any interface.
	Get() V    // Gets the stored "wrapped" value. Useful when you want to work with the value directly. This is the case when nesting wrappers.
}

func New[T WrapperImplementation[V, R], V any, R UnwrapResult]() T {
	// Create a new instance of type T using reflection
	var wrapper T

	// Check if T is a pointer type
	if reflect.TypeOf(wrapper).Kind() == reflect.Ptr {
		// Create a new instance of the type pointed to by T
		v := reflect.New(reflect.TypeOf(wrapper).Elem())
		wrapper = v.Interface().(T) // Type assertion to T
	}

	// Initialize the wrapper
	wrapper.Initialize()

	return wrapper
}

func NewWithValue[T WrapperImplementation[V, R], V any, R UnwrapResult](value V) (*T, error) {
	wrapper := New[T]()

	if err := wrapper.Wrap(value, false); err != nil {
		return nil, err
	}

	return &wrapper, nil
}

func NewWithValueUnsafe[T WrapperImplementation[V, R], V any, R UnwrapResult](value V) T {
	wrapper, err := NewWithValue[T](value)
	if err != nil {
		panic(err)
	}

	return *wrapper
}

// MarshalJSON is a generic implementation of the MarshalJSON method for wrappers. It is used to marshal a wrapper into a JSON value.
// All wrappers should call this method in their MarshalJSON implementation.
func MarshalJSON(wrapper WrapperProvider) ([]byte, error) {
	if reflect.ValueOf(wrapper).IsNil() {
		return nil, fmt.Errorf("marshal into nil wrapper")
	}

	if !wrapper.IsInitialized() { // Make sure the wrapper is initialized since for marshalling we new to automatically initialize the wrapper.
		wrapper.Initialize()
	}

	if wrapper.IsDiscarded() {
		return json.Marshal(nil)
	}

	result, err := json.Marshal(wrapper.UnwrapAny())

	return result, err
}

// UnmarshalJSON is a generic implementation of the UnmarshalJSON method for wrappers. It is used to unmarshal a JSON value into a wrapper.
// All wrappers should call this method in their UnmarshalJSON implementation.
func UnmarshalJSON(data []byte, wrapper WrapperProvider) error {
	if reflect.ValueOf(wrapper).IsNil() {
		return fmt.Errorf("unmarshal into nil wrapper")
	}

	if !wrapper.IsInitialized() { // Same as with marshalling, we need to initialize the wrapper before unmarshalling.
		wrapper.Initialize()
	}

	// Perform regular JSON unmarshalling
	var s interface{}
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	err := wrapper.Wrap(s, false)
	return err
}
