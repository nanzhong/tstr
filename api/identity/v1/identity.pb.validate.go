// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: identity/v1/identity.proto

package identityv1

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on IdentityRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *IdentityRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IdentityRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IdentityRequestMultiError, or nil if none found.
func (m *IdentityRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *IdentityRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return IdentityRequestMultiError(errors)
	}

	return nil
}

// IdentityRequestMultiError is an error wrapping multiple validation errors
// returned by IdentityRequest.ValidateAll() if the designated constraints
// aren't met.
type IdentityRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IdentityRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IdentityRequestMultiError) AllErrors() []error { return m }

// IdentityRequestValidationError is the validation error returned by
// IdentityRequest.Validate if the designated constraints aren't met.
type IdentityRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IdentityRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IdentityRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IdentityRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IdentityRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IdentityRequestValidationError) ErrorName() string { return "IdentityRequestValidationError" }

// Error satisfies the builtin error interface
func (e IdentityRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIdentityRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IdentityRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IdentityRequestValidationError{}

// Validate checks the field values on IdentityResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *IdentityResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on IdentityResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// IdentityResponseMultiError, or nil if none found.
func (m *IdentityResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *IdentityResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetAccessToken()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, IdentityResponseValidationError{
					field:  "AccessToken",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, IdentityResponseValidationError{
					field:  "AccessToken",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetAccessToken()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return IdentityResponseValidationError{
				field:  "AccessToken",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return IdentityResponseMultiError(errors)
	}

	return nil
}

// IdentityResponseMultiError is an error wrapping multiple validation errors
// returned by IdentityResponse.ValidateAll() if the designated constraints
// aren't met.
type IdentityResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m IdentityResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m IdentityResponseMultiError) AllErrors() []error { return m }

// IdentityResponseValidationError is the validation error returned by
// IdentityResponse.Validate if the designated constraints aren't met.
type IdentityResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e IdentityResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e IdentityResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e IdentityResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e IdentityResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e IdentityResponseValidationError) ErrorName() string { return "IdentityResponseValidationError" }

// Error satisfies the builtin error interface
func (e IdentityResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sIdentityResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = IdentityResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = IdentityResponseValidationError{}
