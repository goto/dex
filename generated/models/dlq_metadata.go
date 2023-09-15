// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// DlqMetadata dlq metadata
//
// swagger:model DlqMetadata
type DlqMetadata struct {

	// date
	// Example: 2022-10-24
	Date string `json:"date,omitempty"`

	// size in bytes
	SizeInBytes int64 `json:"size_in_bytes,omitempty"`

	// topic
	Topic string `json:"topic,omitempty"`
}

// Validate validates this dlq metadata
func (m *DlqMetadata) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this dlq metadata based on context it is used
func (m *DlqMetadata) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DlqMetadata) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DlqMetadata) UnmarshalBinary(b []byte) error {
	var res DlqMetadata
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
