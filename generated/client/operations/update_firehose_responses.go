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

// UpdateFirehoseReader is a Reader for the UpdateFirehose structure.
type UpdateFirehoseReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateFirehoseReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateFirehoseOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 400:
		result := NewUpdateFirehoseBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateFirehoseNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewUpdateFirehoseInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUpdateFirehoseOK creates a UpdateFirehoseOK with default headers values
func NewUpdateFirehoseOK() *UpdateFirehoseOK {
	return &UpdateFirehoseOK{}
}

/*
UpdateFirehoseOK describes a response with status code 200, with default header values.

Found firehose with given URN
*/
type UpdateFirehoseOK struct {
	Payload *models.Firehose
}

// IsSuccess returns true when this update firehose o k response has a 2xx status code
func (o *UpdateFirehoseOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this update firehose o k response has a 3xx status code
func (o *UpdateFirehoseOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update firehose o k response has a 4xx status code
func (o *UpdateFirehoseOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this update firehose o k response has a 5xx status code
func (o *UpdateFirehoseOK) IsServerError() bool {
	return false
}

// IsCode returns true when this update firehose o k response a status code equal to that given
func (o *UpdateFirehoseOK) IsCode(code int) bool {
	return code == 200
}

func (o *UpdateFirehoseOK) Error() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseOK  %+v", 200, o.Payload)
}

func (o *UpdateFirehoseOK) String() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseOK  %+v", 200, o.Payload)
}

func (o *UpdateFirehoseOK) GetPayload() *models.Firehose {
	return o.Payload
}

func (o *UpdateFirehoseOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Firehose)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateFirehoseBadRequest creates a UpdateFirehoseBadRequest with default headers values
func NewUpdateFirehoseBadRequest() *UpdateFirehoseBadRequest {
	return &UpdateFirehoseBadRequest{}
}

/*
UpdateFirehoseBadRequest describes a response with status code 400, with default header values.

Update request is not valid.
*/
type UpdateFirehoseBadRequest struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this update firehose bad request response has a 2xx status code
func (o *UpdateFirehoseBadRequest) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update firehose bad request response has a 3xx status code
func (o *UpdateFirehoseBadRequest) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update firehose bad request response has a 4xx status code
func (o *UpdateFirehoseBadRequest) IsClientError() bool {
	return true
}

// IsServerError returns true when this update firehose bad request response has a 5xx status code
func (o *UpdateFirehoseBadRequest) IsServerError() bool {
	return false
}

// IsCode returns true when this update firehose bad request response a status code equal to that given
func (o *UpdateFirehoseBadRequest) IsCode(code int) bool {
	return code == 400
}

func (o *UpdateFirehoseBadRequest) Error() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateFirehoseBadRequest) String() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseBadRequest  %+v", 400, o.Payload)
}

func (o *UpdateFirehoseBadRequest) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateFirehoseBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateFirehoseNotFound creates a UpdateFirehoseNotFound with default headers values
func NewUpdateFirehoseNotFound() *UpdateFirehoseNotFound {
	return &UpdateFirehoseNotFound{}
}

/*
UpdateFirehoseNotFound describes a response with status code 404, with default header values.

Firehose with given URN was not found
*/
type UpdateFirehoseNotFound struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this update firehose not found response has a 2xx status code
func (o *UpdateFirehoseNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update firehose not found response has a 3xx status code
func (o *UpdateFirehoseNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update firehose not found response has a 4xx status code
func (o *UpdateFirehoseNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this update firehose not found response has a 5xx status code
func (o *UpdateFirehoseNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this update firehose not found response a status code equal to that given
func (o *UpdateFirehoseNotFound) IsCode(code int) bool {
	return code == 404
}

func (o *UpdateFirehoseNotFound) Error() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseNotFound  %+v", 404, o.Payload)
}

func (o *UpdateFirehoseNotFound) String() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseNotFound  %+v", 404, o.Payload)
}

func (o *UpdateFirehoseNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateFirehoseNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateFirehoseInternalServerError creates a UpdateFirehoseInternalServerError with default headers values
func NewUpdateFirehoseInternalServerError() *UpdateFirehoseInternalServerError {
	return &UpdateFirehoseInternalServerError{}
}

/*
UpdateFirehoseInternalServerError describes a response with status code 500, with default header values.

internal error
*/
type UpdateFirehoseInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this update firehose internal server error response has a 2xx status code
func (o *UpdateFirehoseInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this update firehose internal server error response has a 3xx status code
func (o *UpdateFirehoseInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this update firehose internal server error response has a 4xx status code
func (o *UpdateFirehoseInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this update firehose internal server error response has a 5xx status code
func (o *UpdateFirehoseInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this update firehose internal server error response a status code equal to that given
func (o *UpdateFirehoseInternalServerError) IsCode(code int) bool {
	return code == 500
}

func (o *UpdateFirehoseInternalServerError) Error() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateFirehoseInternalServerError) String() string {
	return fmt.Sprintf("[PUT /dex/firehoses/{firehoseUrn}][%d] updateFirehoseInternalServerError  %+v", 500, o.Payload)
}

func (o *UpdateFirehoseInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateFirehoseInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
UpdateFirehoseBody update firehose body
swagger:model UpdateFirehoseBody
*/
type UpdateFirehoseBody struct {

	// configs
	Configs *models.FirehoseConfig `json:"configs,omitempty"`

	// description
	// Example: This firehose consumes from booking events and ingests to redis
	Description string `json:"description,omitempty"`
}

// Validate validates this update firehose body
func (o *UpdateFirehoseBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateConfigs(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateFirehoseBody) validateConfigs(formats strfmt.Registry) error {
	if swag.IsZero(o.Configs) { // not required
		return nil
	}

	if o.Configs != nil {
		if err := o.Configs.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "configs")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("body" + "." + "configs")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this update firehose body based on the context it is used
func (o *UpdateFirehoseBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateConfigs(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *UpdateFirehoseBody) contextValidateConfigs(ctx context.Context, formats strfmt.Registry) error {

	if o.Configs != nil {
		if err := o.Configs.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("body" + "." + "configs")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("body" + "." + "configs")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (o *UpdateFirehoseBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *UpdateFirehoseBody) UnmarshalBinary(b []byte) error {
	var res UpdateFirehoseBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
