// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/goto/dex/generated/models"
)

// GetGroupAlertChannelsReader is a Reader for the GetGroupAlertChannels structure.
type GetGroupAlertChannelsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetGroupAlertChannelsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetGroupAlertChannelsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 404:
		result := NewGetGroupAlertChannelsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 500:
		result := NewGetGroupAlertChannelsInternalServerError()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("[GET /dex/subscriptions/groups/{id}/alert_channels] getGroupAlertChannels", response, response.Code())
	}
}

// NewGetGroupAlertChannelsOK creates a GetGroupAlertChannelsOK with default headers values
func NewGetGroupAlertChannelsOK() *GetGroupAlertChannelsOK {
	return &GetGroupAlertChannelsOK{}
}

/*
GetGroupAlertChannelsOK describes a response with status code 200, with default header values.

Successful Operation.
*/
type GetGroupAlertChannelsOK struct {
	Payload *GetGroupAlertChannelsOKBody
}

// IsSuccess returns true when this get group alert channels o k response has a 2xx status code
func (o *GetGroupAlertChannelsOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get group alert channels o k response has a 3xx status code
func (o *GetGroupAlertChannelsOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get group alert channels o k response has a 4xx status code
func (o *GetGroupAlertChannelsOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get group alert channels o k response has a 5xx status code
func (o *GetGroupAlertChannelsOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get group alert channels o k response a status code equal to that given
func (o *GetGroupAlertChannelsOK) IsCode(code int) bool {
	return code == 200
}

// Code gets the status code for the get group alert channels o k response
func (o *GetGroupAlertChannelsOK) Code() int {
	return 200
}

func (o *GetGroupAlertChannelsOK) Error() string {
	return fmt.Sprintf("[GET /dex/subscriptions/groups/{id}/alert_channels][%d] getGroupAlertChannelsOK  %+v", 200, o.Payload)
}

func (o *GetGroupAlertChannelsOK) String() string {
	return fmt.Sprintf("[GET /dex/subscriptions/groups/{id}/alert_channels][%d] getGroupAlertChannelsOK  %+v", 200, o.Payload)
}

func (o *GetGroupAlertChannelsOK) GetPayload() *GetGroupAlertChannelsOKBody {
	return o.Payload
}

func (o *GetGroupAlertChannelsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(GetGroupAlertChannelsOKBody)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetGroupAlertChannelsNotFound creates a GetGroupAlertChannelsNotFound with default headers values
func NewGetGroupAlertChannelsNotFound() *GetGroupAlertChannelsNotFound {
	return &GetGroupAlertChannelsNotFound{}
}

/*
GetGroupAlertChannelsNotFound describes a response with status code 404, with default header values.

Group Not Found Error
*/
type GetGroupAlertChannelsNotFound struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this get group alert channels not found response has a 2xx status code
func (o *GetGroupAlertChannelsNotFound) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get group alert channels not found response has a 3xx status code
func (o *GetGroupAlertChannelsNotFound) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get group alert channels not found response has a 4xx status code
func (o *GetGroupAlertChannelsNotFound) IsClientError() bool {
	return true
}

// IsServerError returns true when this get group alert channels not found response has a 5xx status code
func (o *GetGroupAlertChannelsNotFound) IsServerError() bool {
	return false
}

// IsCode returns true when this get group alert channels not found response a status code equal to that given
func (o *GetGroupAlertChannelsNotFound) IsCode(code int) bool {
	return code == 404
}

// Code gets the status code for the get group alert channels not found response
func (o *GetGroupAlertChannelsNotFound) Code() int {
	return 404
}

func (o *GetGroupAlertChannelsNotFound) Error() string {
	return fmt.Sprintf("[GET /dex/subscriptions/groups/{id}/alert_channels][%d] getGroupAlertChannelsNotFound  %+v", 404, o.Payload)
}

func (o *GetGroupAlertChannelsNotFound) String() string {
	return fmt.Sprintf("[GET /dex/subscriptions/groups/{id}/alert_channels][%d] getGroupAlertChannelsNotFound  %+v", 404, o.Payload)
}

func (o *GetGroupAlertChannelsNotFound) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetGroupAlertChannelsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetGroupAlertChannelsInternalServerError creates a GetGroupAlertChannelsInternalServerError with default headers values
func NewGetGroupAlertChannelsInternalServerError() *GetGroupAlertChannelsInternalServerError {
	return &GetGroupAlertChannelsInternalServerError{}
}

/*
GetGroupAlertChannelsInternalServerError describes a response with status code 500, with default header values.

Internal Error
*/
type GetGroupAlertChannelsInternalServerError struct {
	Payload *models.ErrorResponse
}

// IsSuccess returns true when this get group alert channels internal server error response has a 2xx status code
func (o *GetGroupAlertChannelsInternalServerError) IsSuccess() bool {
	return false
}

// IsRedirect returns true when this get group alert channels internal server error response has a 3xx status code
func (o *GetGroupAlertChannelsInternalServerError) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get group alert channels internal server error response has a 4xx status code
func (o *GetGroupAlertChannelsInternalServerError) IsClientError() bool {
	return false
}

// IsServerError returns true when this get group alert channels internal server error response has a 5xx status code
func (o *GetGroupAlertChannelsInternalServerError) IsServerError() bool {
	return true
}

// IsCode returns true when this get group alert channels internal server error response a status code equal to that given
func (o *GetGroupAlertChannelsInternalServerError) IsCode(code int) bool {
	return code == 500
}

// Code gets the status code for the get group alert channels internal server error response
func (o *GetGroupAlertChannelsInternalServerError) Code() int {
	return 500
}

func (o *GetGroupAlertChannelsInternalServerError) Error() string {
	return fmt.Sprintf("[GET /dex/subscriptions/groups/{id}/alert_channels][%d] getGroupAlertChannelsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetGroupAlertChannelsInternalServerError) String() string {
	return fmt.Sprintf("[GET /dex/subscriptions/groups/{id}/alert_channels][%d] getGroupAlertChannelsInternalServerError  %+v", 500, o.Payload)
}

func (o *GetGroupAlertChannelsInternalServerError) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetGroupAlertChannelsInternalServerError) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

/*
GetGroupAlertChannelsOKBody get group alert channels o k body
swagger:model GetGroupAlertChannelsOKBody
*/
type GetGroupAlertChannelsOKBody struct {

	// alert channels
	AlertChannels []*models.AlertChannel `json:"alert_channels"`
}

// Validate validates this get group alert channels o k body
func (o *GetGroupAlertChannelsOKBody) Validate(formats strfmt.Registry) error {
	var res []error

	if err := o.validateAlertChannels(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetGroupAlertChannelsOKBody) validateAlertChannels(formats strfmt.Registry) error {
	if swag.IsZero(o.AlertChannels) { // not required
		return nil
	}

	for i := 0; i < len(o.AlertChannels); i++ {
		if swag.IsZero(o.AlertChannels[i]) { // not required
			continue
		}

		if o.AlertChannels[i] != nil {
			if err := o.AlertChannels[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGroupAlertChannelsOK" + "." + "alert_channels" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGroupAlertChannelsOK" + "." + "alert_channels" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this get group alert channels o k body based on the context it is used
func (o *GetGroupAlertChannelsOKBody) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := o.contextValidateAlertChannels(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetGroupAlertChannelsOKBody) contextValidateAlertChannels(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(o.AlertChannels); i++ {

		if o.AlertChannels[i] != nil {

			if swag.IsZero(o.AlertChannels[i]) { // not required
				return nil
			}

			if err := o.AlertChannels[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("getGroupAlertChannelsOK" + "." + "alert_channels" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("getGroupAlertChannelsOK" + "." + "alert_channels" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (o *GetGroupAlertChannelsOKBody) MarshalBinary() ([]byte, error) {
	if o == nil {
		return nil, nil
	}
	return swag.WriteJSON(o)
}

// UnmarshalBinary interface implementation
func (o *GetGroupAlertChannelsOKBody) UnmarshalBinary(b []byte) error {
	var res GetGroupAlertChannelsOKBody
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*o = res
	return nil
}
