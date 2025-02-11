// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/goto/dex/generated/models"
)

// UpdateAlertSubscriptionReader is a Reader for the UpdateAlertSubscription structure.
type UpdateAlertSubscriptionReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateAlertSubscriptionReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateAlertSubscriptionOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateAlertSubscriptionBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateAlertSubscriptionNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 409:
		result := NewUpdateAlertSubscriptionConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 422:
		result := NewUpdateAlertSubscriptionUnprocessableEntity()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateAlertSubscriptionInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpdateAlertSubscriptionOK creates a UpdateAlertSubscriptionOK with default headers values
func NewUpdateAlertSubscriptionOK() *UpdateAlertSubscriptionOK {
	return &UpdateAlertSubscriptionOK{}
}

/*
	UpdateAlertSubscriptionOK describes a response with status code 200, with default header values.

Successful Operation.
*/
type UpdateAlertSubscriptionOK struct {
	Payload *UpdateAlertSubscriptionOKBody
}

func (o *UpdateAlertSubscriptionOK) Error() string {
	return fmt.Sprintf("[PUT /dex/subscriptions/{id}][%d] updateAlertSubscriptionOK  %+v", 200, o.Payload)
}
func (o *UpdateAlertSubscriptionOK) GetPayload() *UpdateAlertSubscriptionOKBody {
	return o.Payload
}

func (o *UpdateAlertSubscriptionOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(UpdateAlertSubscriptionOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAlertSubscriptionBadRequest creates a UpdateAlertSubscriptionBadRequest with default headers values
func NewUpdateAlertSubscriptionBadRequest() *UpdateAlertSubscriptionBadRequest {
	return &UpdateAlertSubscriptionBadRequest{}
}

/*
	UpdateAlertSubscriptionBadRequest describes a response with status code 400, with default header values.

Validation Error
*/
type UpdateAlertSubscriptionBadRequest struct {
	Payload *models.ErrorResponse
}

func (o *UpdateAlertSubscriptionBadRequest) Error() string {
	return fmt.Sprintf("[PUT /dex/subscriptions/{id}][%d] updateAlertSubscriptionBadRequest  %+v", 400, o.Payload)
}
func (o *UpdateAlertSubscriptionBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateAlertSubscriptionBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAlertSubscriptionNotFound creates a UpdateAlertSubscriptionNotFound with default headers values
func NewUpdateAlertSubscriptionNotFound() *UpdateAlertSubscriptionNotFound {
	return &UpdateAlertSubscriptionNotFound{}
}

/*
	UpdateAlertSubscriptionNotFound describes a response with status code 404, with default header values.

Not Found Error
*/
type UpdateAlertSubscriptionNotFound struct {
	Payload *models.ErrorResponse
}

func (o *UpdateAlertSubscriptionNotFound) Error() string {
	return fmt.Sprintf("[PUT /dex/subscriptions/{id}][%d] updateAlertSubscriptionNotFound  %+v", 404, o.Payload)
}
func (o *UpdateAlertSubscriptionNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateAlertSubscriptionNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAlertSubscriptionConflict creates a UpdateAlertSubscriptionConflict with default headers values
func NewUpdateAlertSubscriptionConflict() *UpdateAlertSubscriptionConflict {
	return &UpdateAlertSubscriptionConflict{}
}

/*
	UpdateAlertSubscriptionConflict describes a response with status code 409, with default header values.

Duplicate subscription
*/
type UpdateAlertSubscriptionConflict struct {
	Payload *models.ErrorResponse
}

func (o *UpdateAlertSubscriptionConflict) Error() string {
	return fmt.Sprintf("[PUT /dex/subscriptions/{id}][%d] updateAlertSubscriptionConflict  %+v", 409, o.Payload)
}
func (o *UpdateAlertSubscriptionConflict) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateAlertSubscriptionConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAlertSubscriptionUnprocessableEntity creates a UpdateAlertSubscriptionUnprocessableEntity with default headers values
func NewUpdateAlertSubscriptionUnprocessableEntity() *UpdateAlertSubscriptionUnprocessableEntity {
	return &UpdateAlertSubscriptionUnprocessableEntity{}
}

/*
	UpdateAlertSubscriptionUnprocessableEntity describes a response with status code 422, with default header values.

Missing namespace and slack channel in shield
*/
type UpdateAlertSubscriptionUnprocessableEntity struct {
	Payload *models.ErrorResponse
}

func (o *UpdateAlertSubscriptionUnprocessableEntity) Error() string {
	return fmt.Sprintf("[PUT /dex/subscriptions/{id}][%d] updateAlertSubscriptionUnprocessableEntity  %+v", 422, o.Payload)
}
func (o *UpdateAlertSubscriptionUnprocessableEntity) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateAlertSubscriptionUnprocessableEntity) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateAlertSubscriptionInternalServerError creates a UpdateAlertSubscriptionInternalServerError with default headers values
func NewUpdateAlertSubscriptionInternalServerError() *UpdateAlertSubscriptionInternalServerError {
	return &UpdateAlertSubscriptionInternalServerError{}
}

/*
	UpdateAlertSubscriptionInternalServerError describes a response with status code 500, with default header values.

Internal Error
*/
type UpdateAlertSubscriptionInternalServerError struct {
	Payload *models.ErrorResponse
}

func (o *UpdateAlertSubscriptionInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /dex/subscriptions/{id}][%d] updateAlertSubscriptionInternalServerError  %+v", 500, o.Payload)
}
func (o *UpdateAlertSubscriptionInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateAlertSubscriptionInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
UpdateAlertSubscriptionOKBody update alert subscription o k body
swagger:model UpdateAlertSubscriptionOKBody
*/
type UpdateAlertSubscriptionOKBody struct {

	// subscription
	Subscription *models.Subscription `json:"subscription,omitempty"`
}

// Validate validates this update alert subscription o k body
func (o *UpdateAlertSubscriptionOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateSubscription(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateAlertSubscriptionOKBody) validateSubscription(formats strfmt.Registry) error {
	if swag.IsZero(o.Subscription) { // not required
		return nil
	}

	if o.Subscription != nil {
		if err := o.Subscription.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("updateAlertSubscriptionOK" + "." + "subscription")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this update alert subscription o k body based on the context it is used
func (o *UpdateAlertSubscriptionOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateSubscription(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateAlertSubscriptionOKBody) contextValidateSubscription(ctx context.Context, formats strfmt.Registry) error {

	if o.Subscription != nil {
		if err := o.Subscription.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("updateAlertSubscriptionOK" + "." + "subscription")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateAlertSubscriptionOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateAlertSubscriptionOKBody) UnmarshalBinary(b []byte) error {
	var res UpdateAlertSubscriptionOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
