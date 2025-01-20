# Wrappers

![Go](https://img.shields.io/badge/Go-1.23.4%2B-blue.svg)
![License](https://img.shields.io/badge/License-MIT-green.svg)
[![Build Go application](https://github.com/zealsprince/wrappers/actions/workflows/go-build.yml/badge.svg)](https://github.com/zealsprince/wrappers/actions/workflows/go-build.yml)
[![Test Go application](https://github.com/zealsprince/wrappers/actions/workflows/go-test.yml/badge.svg)](https://github.com/zealsprince/wrappers/actions/workflows/go-test.yml)

Wrappers is a versatile and extensible library for Go that provides a comprehensive validation layer for struct fields. By leveraging generic wrappers and regex-based validations, it ensures type safety, data integrity, and streamlined JSON marshalling/unmarshalling.

**For motivations on this package, refer to the [Motivations](#motivations) section.**

## Features

- **Generic Wrappers**: Easily wrap and validate any data type with built-in error handling.
- **Automatic Type Inference**: Wraps `any` and always unwraps to an elementary type of `bool | int64 | float64 | string`, ideal for data serialization tasks.
- **Regex-Based Validations**: Reduce boilerplate by using reusable regex-based wrappers for common validation needs like emails, phone numbers, URLs, etc.
- **Seamless JSON Integration**: Automatic validation during JSON marshalling and unmarshalling with optional discarding (defaulting) of values.
- **Extensible Architecture**: Easily extend the library with custom validation logic as needed.
- **Performance Optimizations**: Common validations are cache compiled regex patterns to enhance performance.

## What are Wrappers?

Essentially, you'll be used to something like this:

```go
type User struct {
    ID    int64  `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Unmarshal JSON data into a User struct.
    var user User
    err := json.Unmarshal([]byte(`{"id": "1", "name": "Alice", "email": "example@test.com"}`), &user)
    if err != nil {
        panic(err)
    }

    // Validate the user's email address.
    if !isValidEmail(user.Email) {
        return errors.New("Invalid email address")
    }

    // Validate the user's ID. It comes as a string, but our code will work with integers.
    id, err := strconv.ParseInt(user.ID, 10, 64)
    if err != nil {
        return errors.New("Invalid user ID")
    } else if id <= 0 {
        return errors.New("Invalid user ID")
    }

    // Validate the user's name. Can't be empty!
    if user.Name == "" {
        return errors.New("Invalid user name")
    }

    // Continue with the rest of the application logic...
}
```

But wouldn't it be nice to have something like this instead?

```go
type User struct {
    ID    *custom.WrapperPositiveInt `json:"id"`
    Name  *wrappers.WrapperString    `json:"name"`
    Email *regex.WrapperRegexEmail   `json:"email"`
}

func main() {
    // Unmarshal JSON data into a User struct.
    var user User
    err := json.Unmarshal([]byte(`{"id": 1, "name": "Alice", "email": "example@test.com"}`), &user)
    if err != nil {
        panic(err) // If there was a failing validation, this would be an error!
    }

    // And not have to do anything else; it's all handled for you!

    // Plus, if you want to marshal it back to JSON, it will be the same as the input!
```

You can chose to discard values without raising errors

## Overview

The wrappers package provides a structured approach to manage and validate various data types in Go through encapsulation in dedicated wrapper types.

### On their own

**Wrapping**: Each wrapper type (e.g., WrapperInt, WrapperString) allows elementary values to be stored safely as any implementing type. The Wrap method performs validation on the input, ensuring that only acceptable values are assigned. If the input is invalid, the wrapper can either return an error or discard the value based on the method parameters.

**Discarding**: This process happens whenever the wrapped data fails the validation. A discard operation always flags a wrapper as discarded which further processes can refer to. It then nullifies its value to the default value for that type. Normally, this also throws an error during Wrapping or Unmarshalling processes.

>[!TIP]
> Wrapping takes an optional `discard` boolean parameter which will cause errors to be dropped but still discard the value. For unmarshalling processes, the `Discarder` type can be used to suppress errors and discard values silently.

**Unwrapping**: The Unwrap method retrieves the stored value from the wrapper. If the value is marked as discarded, a default value (e.g., 0 for integers) is returned, helping to maintain clarity about the data's state.

### Within data structures

**Marshalling**: The MarshalJSON method enables wrappers to be serialized into JSON format. When a struct containing wrappers is marshalled, discarded values are serialized as null or omitted (if marked with omitempty and nilled), ensuring accurate representation.

**Unmarshalling**: The UnmarshalJSON method allows JSON data to populate the corresponding wrappers. It validates the incoming data, storing values if they fit the expected format or discarding them otherwise, while also handling errors appropriately.

> [!WARNING]
> This requires further handling to ensure no erroneous data is processed. It is generally recommended - especially during unmarshalling - to handle the unmarshal errors and not proceed with invalid data.

## Sub-packages

Wrappers ships with two additional sub-packages that enable further every day usage:

- Regex: This package contains the `WrapperRegex` which inherits the functionality of WrapperString but extends it with automated Regex validation. The package additionally ships with a core set of common validations.
- Enum: This package contains the `WrapperEnum` which allows for handling of custom single value data types and the additional checks for conformity to the core enumerating type.

## Usage

### Standalone Wrapper Usage

The library provides a generic `Wrapper` foundational type that is implemented for various sub-types. By default the package ships with elementary type wrappers as well as a few others such as `time.Time`.

```go
package main

import (
    "fmt"
    "wrappers"
)

func main() {
    // A simple string wrapper with direct initialization.
    stringWrapper, err := wrappers.NewWithValue[*wrappers.WrapperString]("Hello, World!")
    if err != nil {
        panic(err)
    }

    // Unwrap the string value.
    fmt.Println("Unwrapped string value:", stringWrapper.Unwrap())
    // (string) "Hello, World!"

    // A simple int wrapper with safe wrapping of an originally string value.
    intWrapper := wrappers.New[*wrappers.WrapperInt]()
    err = intWrapper.Wrap("1", false)
    if err != nil {
        panic(err)
    }

    // Unwrap the integer value.
    fmt.Println("Unwrapped int value:", intWrapper.Unwrap())
    // (int64) 1
    
    // Due to the nature of wrappers, we can pass one into the other.
    err = stringWrapper.Wrap(stringWrapper, false)
    if err != nil {
        panic(err)
    }

    // Unwrap the string value which was wrapped from the intWrapper.
    fmt.Println("Unwrapped string value:", stringWrapper.Unwrap())
    // (string) "1"
}
```

Notice that handling and type conversion is performed by the wrappers. However, during cases where an invalid value is passed, wrappers can discard the value without raising an error. This should generally be avoided or used in strict combination with the `IsDiscarded` method to assure no erroneous is further handled and causing side-effects.

```go
    // [...] Continuation of previous code

    // Let's wrap an invalid value which would result in an error.
    // To avoid handling the error, we can pass the discard flag as true which
    // will discard the value and not return an error.
    intWrapper.Wrap("Hello, World!", true)

    // Let's take a look at the value now!
    fmt.Println("Unwrapped int value:", intWrapper.Unwrap())
    // (int64) 0

    // Notice the value was set to the int64 default.
    // It is referred to as discarded from here on out.
    fmt.Println("Is discarded?", intWrapper.IsDiscarded())
    // true
```

### Struct Integration

Integrate wrappers directly into your structs to ensure that validations occur automatically during JSON operations and data assignments.

```go
package main

import (
    "encoding/json"
    "fmt"
    "wrappers/regex"
)

type Data struct {
    Name  *wrappers.WrapperString  `json:"name"`
    Email *regex.WrapperRegexEmail `json:"email,omitempty"`
}

func main() {
    // Assuming the given payload.
    input := "{\"name\":\"Andrew\",\"email\":\"andrew@example.com\"}"

    // Let's unmarshal into a data struct.
    var data Data
    err := json.Unmarshal([]byte(input), &data)
    if err != nil {
        fmt.Println("Error unmarshalling JSON:", err)
        return
    }

    // Let's print the unwrapped email.
    fmt.Println("Email:", data.Email.Unwrap())

    // Let's marshal the data struct back to JSON.
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Error marshalling to JSON:", err)
        return
    }

    // Marshalling the struct back to JSON will result in the same data:
    // {"name": "Andrew", "email": "andrew@example.com"}
}
```

#### Omitting Data

Just like with normal structs, adding `,omitempty` will remove the value from the resulting marshalled struct. However, to accomplish this effectively and nullify data, you will have to nil out the value in your nesting struct.

```go
    // [...] Continuation of previous code

    data.Email = nil

    // Let's marshal the data struct back to JSON, again.
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Println("Error marshalling to JSON:", err)
        return
    }

    // Marshalling this time will result in a smaller version of the data:
    // {"name": "Andrew"}
```

#### Handling Faulty Data during Unmarshalling

The Discarder type is used to suppress errors during JSON unmarshalling, allowing for automatic handling of invalid values without raising errors.

Keep in mind, this scenario assumes that you will have to check if the value was discarded or not before proceeding with further operations.

**It is safest to simply handle the Unmarshal error and not proceed with the data if it is invalid.**

```go
// [...] package and imports

type DiscardingData struct {
    Value *wrappers.Discarder[*wrappers.WrapperString] `json:"value"`
}

func main() {
    // Example with valid string input
    input := `{"value": "Hello, World!"}`
    data := DiscardingData{}
    
    // Due to using a Discarder, we won't receive any errors.
    json.Unmarshal([]byte(input), &data)
    
    fmt.Println("Unwrapped value:", nested.Value.Proxy.Unwrap())
    // (string) "Hello, World!"

    // Example with invalid input that will be discarded
    invalidData := `{"value": {"key": "value"}}`
    
    // Same in this case, except this time the result is discarded.
    json.Unmarshal([]byte(invalidData), &nested)
    
    // Your further code can now check if the value was discarded.
    fmt.Printf("Discarded: %v / Unwrapped result: %v\n", nested.Value.Proxy.IsDiscarded(), nested.Value.Proxy.Unwrap())
    // Discarded: true / Unwrapped result: ""
}
```

## Creating Custom Regex Wrappers

While the `regex` sub-package covers many common validation scenarios, you can create custom wrappers tailored to your specific needs by following these steps:

1. **Define the Wrapper Struct**

    Embed the generic `Wrapper` and include any additional fields or methods as necessary.

    ```go
    package regex

    import (
        "wrappers"
    )

    type WrapperRegexCustom struct {
        WrapperRegex
    }
    ```

2. **Implement Initialization Method**

    Initialize the embedded `WrapperRegex` with your custom regex pattern.

    ```go
    func (wrapper *WrapperRegexCustom) Initialize() {
        wrapper.WrapperRegex.SetPattern(
            "WrapperRegexCustom",
            "your-regex-pattern"
        )
        wrapper.WrapperBase.Initialize() // Make sure to call the base initialization method.
    }
    ```

3. **Proxy the UnmarshalJSON with Initialization**

    The only thing left is to validate that your custom Wrapper initializes its embedding WrapperRegex in cases of instantiation during Unmarshalling.

    ```go
    func (wrapper *WrapperRegexCustom) UnmarshalJSON(data []byte) error {
        if !wrapper.IsInitialized() {
            wrapper.Initialize()
        }
        return wrapper.WrapperRegex.UnmarshalJSON(data)
    }
    ```

This approach ensures that your custom wrappers are consistent with existing ones, leveraging the underlying validation logic provided by `WrapperRegex`.

### Typed Wrappers

To implement a completely new typed wrapper, please refer to the existing implementations. A good example would be the `WrapperCountry` type within the root wrapper package.

## Motivations

I conceived the idea of this library while working at [Savages Corp](https://github.com/savages-corp) building our [Data Layer](https://data-layer.com/) project. One of my daily activities while writing integrations was matching API data structures and having to validate each and every field after unmarshalling to structs. A problem that kept repeating itself is data validation. Generally, you will unmarshal to a struct and then have to step through the fields to make sure everything is fine and no one on the integration side (especially in direct customer managed environments) has inevitably changed a field's type, format or structure. This leads to a cumbersome cat-and-mouse game of constantly catching up and fixing bugs time and time again.

**What also follows is a massive amount of boilerplate validation code.**

Wrappers was created to address several pervasive challenges developers encounter when managing and validating data within Go applications. Its creation was driven by the need for a more efficient, type-safe, and extensible approach to handling diverse data types, especially in contexts involving JSON serialization and deserialization. Below are the primary motivations behind developing Wrappers:

1. Enhanced Data Validation: In many Go applications, ensuring the integrity and validity of data—whether coming from user input, external APIs, or databases—is paramount. Traditional approaches often involve repetitive boilerplate code to perform type assertions, range checks, and format validations. Wrappers streamline this process by encapsulating validation logic within dedicated types, reducing redundancy and minimizing the likelihood of human error.

2. Type Safety and Generics Leveraging: Go's type system is robust, but when dealing with generic data structures or interfaces like any, maintaining type safety can become cumbersome. Wrappers harness Go's generics to provide a type-safe mechanism for wrapping and unwrapping values. This ensures that data transformations are explicit and safeguarded against type mismatches, enhancing overall code reliability.

3. Seamless JSON Integration: JSON is a ubiquitous data interchange format, and Go's encoding/json package is widely used for serialization and deserialization. However, integrating complex validation logic directly within structs can lead to verbose and hard-to-maintain code. Wrappers automate validation during JSON operations, ensuring that data adheres to specified formats and constraints without cluttering the business logic with repetitive checks.

4. Reduction of Boilerplate Code: Manual validation and type conversion often result in repetitive code patterns that are both time-consuming to write and difficult to maintain. By providing generic wrappers and reusable validation mechanisms—such as regex-based validators—Wrappers significantly reduce the need for boilerplate code. This allows developers to focus on core application logic rather than mundane validation tasks.

5. Extensibility for Custom Validation Needs: Every application has unique data validation requirements. Wrappers are designed with extensibility in mind, allowing developers to create custom wrappers tailored to specific validation rules or data formats. This flexibility ensures that Wrappers can adapt to a wide range of use cases, from simple type checks to complex pattern validations.

6. Improved Error Handling and Data Integrity: Handling invalid data gracefully is crucial for building resilient applications. Wrappers incorporate sophisticated error handling mechanisms, such as the Discarder type, which allows developers to control whether invalid data should result in errors or be silently discarded. This nuanced control helps maintain data integrity while providing flexibility in how applications respond to erroneous inputs.

Ultimately, Wrappers aims to simplify repeating development and data handling patterns. It takes a bit of getting used to and has some nuances and code style requirements, but it is a powerful tool for managing data validation *while also retaining a large portion of your sanity*.

## License

This project is licensed under the [MIT License](LICENSE).
