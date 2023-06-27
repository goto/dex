// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// SubscriptionForm subscription form
//
// swagger:model SubscriptionForm
type SubscriptionForm struct {

	// alert severity
	// Required: true
	// Enum: [INFO WARNING CRITICAL]
	AlertSeverity *string `json:"alert_severity"`

	// channel criticality
	// Required: true
	// Enum: [INFO WARNING CRITICAL]
	ChannelCriticality *string `json:"channel_criticality"`

	// Shield's group id
	// Example: 913464a1-4f87-4312-b4c8-5ac3ebd4f1f2
	// Required: true
	GroupID *string `json:"group_id"`

	// project id
	// Example: p-gojek-id
	// Required: true
	ProjectID *string `json:"project_id"`

	// Firehose Unique Identifier
	// Required: true
	ResourceID *string `json:"resource_id"`

	// resource type
	// Required: true
	// Enum: [firehose]
	ResourceType *string `json:"resource_type"`
}

// Validate validates this subscription form
func (m *SubscriptionForm) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAlertSeverity(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateChannelCriticality(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateGroupID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProjectID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResourceType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

var subscriptionFormTypeAlertSeverityPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["INFO","WARNING","CRITICAL"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		subscriptionFormTypeAlertSeverityPropEnum = append(subscriptionFormTypeAlertSeverityPropEnum, v)
	}
}

const (

	// SubscriptionFormAlertSeverityINFO captures enum value "INFO"
	SubscriptionFormAlertSeverityINFO string = "INFO"

	// SubscriptionFormAlertSeverityWARNING captures enum value "WARNING"
	SubscriptionFormAlertSeverityWARNING string = "WARNING"

	// SubscriptionFormAlertSeverityCRITICAL captures enum value "CRITICAL"
	SubscriptionFormAlertSeverityCRITICAL string = "CRITICAL"
)

// prop value enum
func (m *SubscriptionForm) validateAlertSeverityEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, subscriptionFormTypeAlertSeverityPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SubscriptionForm) validateAlertSeverity(formats strfmt.Registry) error {

	if err := validate.Required("alert_severity", "body", m.AlertSeverity); err != nil {
		return err
	}

	// value enum
	if err := m.validateAlertSeverityEnum("alert_severity", "body", *m.AlertSeverity); err != nil {
		return err
	}

	return nil
}

var subscriptionFormTypeChannelCriticalityPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["INFO","WARNING","CRITICAL"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		subscriptionFormTypeChannelCriticalityPropEnum = append(subscriptionFormTypeChannelCriticalityPropEnum, v)
	}
}

const (

	// SubscriptionFormChannelCriticalityINFO captures enum value "INFO"
	SubscriptionFormChannelCriticalityINFO string = "INFO"

	// SubscriptionFormChannelCriticalityWARNING captures enum value "WARNING"
	SubscriptionFormChannelCriticalityWARNING string = "WARNING"

	// SubscriptionFormChannelCriticalityCRITICAL captures enum value "CRITICAL"
	SubscriptionFormChannelCriticalityCRITICAL string = "CRITICAL"
)

// prop value enum
func (m *SubscriptionForm) validateChannelCriticalityEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, subscriptionFormTypeChannelCriticalityPropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SubscriptionForm) validateChannelCriticality(formats strfmt.Registry) error {

	if err := validate.Required("channel_criticality", "body", m.ChannelCriticality); err != nil {
		return err
	}

	// value enum
	if err := m.validateChannelCriticalityEnum("channel_criticality", "body", *m.ChannelCriticality); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionForm) validateGroupID(formats strfmt.Registry) error {

	if err := validate.Required("group_id", "body", m.GroupID); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionForm) validateProjectID(formats strfmt.Registry) error {

	if err := validate.Required("project_id", "body", m.ProjectID); err != nil {
		return err
	}

	return nil
}

func (m *SubscriptionForm) validateResourceID(formats strfmt.Registry) error {

	if err := validate.Required("resource_id", "body", m.ResourceID); err != nil {
		return err
	}

	return nil
}

var subscriptionFormTypeResourceTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["firehose"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		subscriptionFormTypeResourceTypePropEnum = append(subscriptionFormTypeResourceTypePropEnum, v)
	}
}

const (

	// SubscriptionFormResourceTypeFirehose captures enum value "firehose"
	SubscriptionFormResourceTypeFirehose string = "firehose"
)

// prop value enum
func (m *SubscriptionForm) validateResourceTypeEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, subscriptionFormTypeResourceTypePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *SubscriptionForm) validateResourceType(formats strfmt.Registry) error {

	if err := validate.Required("resource_type", "body", m.ResourceType); err != nil {
		return err
	}

	// value enum
	if err := m.validateResourceTypeEnum("resource_type", "body", *m.ResourceType); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this subscription form based on context it is used
func (m *SubscriptionForm) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SubscriptionForm) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriptionForm) UnmarshalBinary(b []byte) error {
	var res SubscriptionForm
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
